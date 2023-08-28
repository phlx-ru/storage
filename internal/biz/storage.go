package biz

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gosimple/slug"
	"github.com/phlx-ru/hatchet/logger"
	"github.com/phlx-ru/hatchet/metrics"

	v1 "storage/api/storage/v1"
	"storage/ent"
	"storage/internal/clients/auth"
	"storage/internal/clients/minio"
	"storage/internal/conf"
)

const (
	metricPrefix = `biz.storage`

	defaultContentType = `application/octet-stream`
)

type StorageUsecase struct {
	authClient                  auth.Client
	minioClient                 minio.Client
	fileRepo                    fileRepository
	auth                        *conf.Auth
	metric                      metrics.Metrics
	logger                      *log.Helper
	useAuthorizationForDownload bool
}

func NewStorageUsecase(
	authClient auth.Client,
	minioClient minio.Client,
	fileRepo fileRepository,
	auth *conf.Auth,
	metric metrics.Metrics,
	logs log.Logger,
) *StorageUsecase {
	loggerHelper := logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", metricPrefix)
	return &StorageUsecase{
		authClient:                  authClient,
		minioClient:                 minioClient,
		fileRepo:                    fileRepo,
		auth:                        auth,
		metric:                      metric,
		logger:                      loggerHelper,
		useAuthorizationForDownload: false,
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

func makeObjectPath(userID int, filename string) string {
	return fmt.Sprintf(`%d/%s`, userID, slugFromFilename(filename))
}

func (s *StorageUsecase) Upload(ctx context.Context, file *UploadFile) (*ent.File, error) {
	if file == nil {
		return nil, fmt.Errorf(`file is empty`)
	}
	var err error
	var user *AuthenticatedUser
	if !s.isIntegrations(ctx) {
		if user, err = s.user(ctx); err != nil {
			return nil, err
		}
	}

	userID := 0 // for integrations
	if user != nil {
		userID = int(user.ID())
	}

	contentType := mime.TypeByExtension(filepath.Ext(file.Filename))
	if contentType == "" {
		contentType = defaultContentType
	}

	objectPath := makeObjectPath(userID, file.Filename)
	found, err := s.fileRepo.FindByObjectPath(ctx, objectPath)
	if err != nil && ent.IsNotFound(err) {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	if found != nil {
		return nil, v1.ErrorValidationFailed(`file with object path [%s] is already exists`, objectPath)
	}

	saved, err := s.fileRepo.Create(ctx, &ent.File{
		UserID:     userID,
		Filename:   file.Filename,
		ObjectPath: objectPath,
		Size:       int(file.Size),
		MimeType:   contentType,
		DeletedAt:  pointer.ToTime(time.Now()), // For restore after success upload
	})
	if err != nil {
		return nil, err
	}

	_, err = s.minioClient.UploadFromReader(ctx, file.Reader, file.Size, contentType, objectPath)

	if err == nil {
		err = s.fileRepo.Restore(ctx, saved.UID.String())
	}

	return saved, err
}

func (s *StorageUsecase) Download(ctx context.Context, uid string, writer gin.ResponseWriter) error {
	if s.useAuthorizationForDownload && !s.isIntegrations(ctx) {
		if _, err := s.user(ctx); err != nil {
			return err
		}
	}

	f, err := s.fileRepo.FindByUID(ctx, uid)
	if err != nil {
		return err
	}
	writer.Header().Set(`Content-Type`, f.MimeType)
	disposition := "inline"
	if !strings.HasPrefix(f.MimeType, "image/") {
		disposition = fmt.Sprintf(`attachment; filename="%s"`, f.Filename)
	}
	writer.Header().Set(`Content-Disposition`, disposition)
	return s.minioClient.DownloadToWriter(ctx, writer, f.ObjectPath)
}

func (s *StorageUsecase) FilesList(ctx context.Context) ([]*ent.File, error) {
	var err error
	var user *AuthenticatedUser
	if !s.isIntegrations(ctx) {
		if user, err = s.user(ctx); err != nil {
			return nil, err
		}
	}

	userID := 0 // for integrations
	if user != nil {
		userID = int(user.ID())
	}

	limit := 100
	offset := 0
	files, err := s.fileRepo.FindByUserID(ctx, userID, limit, offset)

	return files, err
}
