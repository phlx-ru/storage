package biz

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gosimple/slug"

	"storage/ent"
	"storage/internal/clients/auth"
	"storage/internal/clients/yandex"
	"storage/internal/pkg/logger"
	"storage/internal/pkg/metrics"
	"storage/internal/pkg/strings"
)

const (
	metricPrefix = `biz.storage`
)

type StorageUsecase struct {
	authClient   auth.Client
	yandexClient yandex.Client
	fileRepo     FileRepo
	metric       metrics.Metrics
	logger       *log.Helper
}

func NewStorageUsecase(
	authClient auth.Client,
	yandexClient yandex.Client,
	fileRepo FileRepo,
	metric metrics.Metrics,
	logs log.Logger,
) *StorageUsecase {
	return &StorageUsecase{
		authClient:   authClient,
		yandexClient: yandexClient,
		fileRepo:     fileRepo,
		metric:       metric,
		logger:       logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "biz/storage"),
	}
}

type UploadFile struct {
	Reader   io.Reader
	Size     int64
	Filename string
}

func slugFromFilename(filename string) string {
	ext := filepath.Ext(filename)
	return slug.Make(filename[0:len(filename)-len(ext)]) + ext
}

func makeObjectPath(userID int64, filename string) string {
	return fmt.Sprintf(`%d/%s`, userID, slugFromFilename(filename))
}

func (s *StorageUsecase) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		s.logger.WithContext(ctx).Errorf(`biz storage method %s failed: %v`, method, err)
		s.metric.Increment(strings.Metric(metricPrefix, method, `failure`))
	} else {
		s.metric.Increment(strings.Metric(metricPrefix, method, `success`))
	}
}

func (s *StorageUsecase) Upload(ctx context.Context, file *UploadFile, authToken string) (*ent.File, error) {
	method := `upload`
	defer s.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { s.postProcess(ctx, method, err) }()

	check, err := s.authClient.Check(ctx, authToken)
	if err != nil {
		return nil, err
	}
	contentType := mime.TypeByExtension(filepath.Ext(file.Filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	objectPath := makeObjectPath(check.User.ID, file.Filename)

	saved, err := s.fileRepo.Create(ctx, &ent.File{
		UserID:     int(check.User.ID),
		Filename:   file.Filename,
		ObjectPath: objectPath,
		Size:       int(file.Size),
		MimeType:   contentType,
		DeletedAt:  pointer.ToTime(time.Now()),
	})

	_, err = s.yandexClient.UploadFromReader(ctx, file.Reader, file.Size, contentType, objectPath)

	if err == nil {
		err = s.fileRepo.Restore(ctx, saved.UID.String())
	}

	return saved, err
}

func (s *StorageUsecase) Download(ctx context.Context, uid string, writer io.Writer) error {
	method := `download`
	defer s.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { s.postProcess(ctx, method, err) }()

	file, err := s.fileRepo.FindByUID(ctx, uid)
	if err != nil {
		return err
	}

	err = s.yandexClient.DownloadToWriter(ctx, writer, file.ObjectPath)

	return err
}

func (s *StorageUsecase) FilesList(ctx context.Context, token string) ([]*ent.File, error) {
	method := `filesList`
	defer s.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { s.postProcess(ctx, method, err) }()

	check, err := s.authClient.Check(ctx, token)
	if err != nil {
		return nil, err
	}

	limit := 100
	offset := 0
	files, err := s.fileRepo.FindByUserID(ctx, int(check.User.ID), limit, offset)

	return files, err
}
