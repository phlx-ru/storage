package service

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/phlx-ru/hatchet/logger"
	"github.com/phlx-ru/hatchet/metrics"
	"github.com/phlx-ru/hatchet/watcher"

	v1 "storage/api/storage/v1"
	"storage/internal/biz"
	storage "storage/schema"
	storageComponents "storage/schema/storage"
)

const (
	metricPrefix = `service.storage`
	uidMaxLength = 36
)

type StorageService struct {
	usecase *biz.StorageUsecase
	metric  metrics.Metrics
	logger  *log.Helper
	watcher *watcher.Watcher
}

func NewGatewayService(
	usecase *biz.StorageUsecase,
	metric metrics.Metrics,
	logs log.Logger,
) *StorageService {
	loggerHelper := logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", metricPrefix)
	return &StorageService{
		usecase: usecase,
		metric:  metric,
		logger:  loggerHelper,
		watcher: watcher.New(metricPrefix, loggerHelper, metric).
			WithWarningErrorChecks([]func(err error) bool{
				v1.IsAccessDenied,
				v1.IsValidationFailed,
				v1.IsUnauthorized,
			}),
	}
}

func (s *StorageService) GetSwagger(c *gin.Context) {
	var err error
	defer s.watcher.OnPreparedMethod(`GetSwagger`).Results(func() (context.Context, error) {
		return c.Request.Context(), err
	})

	swagger, err := storage.GetSwagger()
	if err != nil {
		s.responseError(c, errors.InternalServer(`internal_error`, err.Error()))
		return
	}
	swagger.InternalizeRefs(c.Request.Context(), nil) // TODO Bug with nested refs and properties includes
	s.responseOK(c, swagger)
}

func (s *StorageService) Download(c *gin.Context, uid storageComponents.Uid) {
	var err error
	defer s.watcher.OnPreparedMethod(`Download`).Results(func() (context.Context, error) {
		return c.Request.Context(), err
	})

	if uid == "" {
		s.responseError(c, fmt.Errorf(`UID is empty`))
		return
	}

	if len(uid) > uidMaxLength {
		s.responseError(c, fmt.Errorf(`UID must have 36 symbols or less`))
		return
	}

	err = s.usecase.Download(c.Request.Context(), uid, c.Writer)
	if err != nil {
		s.responseError(c, err)
	}
}

func (s *StorageService) DownloadOptions(c *gin.Context, _ storageComponents.Uid) {
	var err error
	defer s.watcher.OnPreparedMethod(`DownloadOptions`).Results(func() (context.Context, error) {
		return c.Request.Context(), err
	})

	s.responseNoContent(c)
}

func (s *StorageService) Upload(c *gin.Context, params storage.UploadParams) {
	var err error
	defer s.watcher.OnPreparedMethod(`Upload`).Results(func() (context.Context, error) {
		return c.Request.Context(), err
	})

	if err = s.validate(params); err != nil {
		s.responseValidationError(c, err)
		return
	}

	uploadFile := &biz.UploadFile{
		Reader:   c.Request.Body,
		Size:     c.Request.ContentLength,
		Filename: params.Filename,
	}

	file, err := s.usecase.Upload(c.Request.Context(), uploadFile)
	if err != nil {
		s.responseError(c, err)
		return
	}

	response := &storageComponents.UploadResponse{
		Filename:   file.Filename,
		MimeType:   pointer.ToString(file.MimeType),
		ObjectPath: file.ObjectPath,
		Size:       pointer.ToInt(file.Size),
		Uid:        file.UID.String(),
		UserId:     file.UserID,
	}

	s.responseOK(c, response)
}

func (s *StorageService) FilesList(c *gin.Context) {
	var err error
	defer s.watcher.OnPreparedMethod(`FilesList`).Results(func() (context.Context, error) {
		return c.Request.Context(), err
	})

	files, err := s.usecase.FilesList(c.Request.Context())
	if err != nil {
		s.responseError(c, errors.InternalServer(`internal_server`, err.Error()))
		return
	}

	filesList := make([]storageComponents.FileItemCompact, len(files))
	for _, file := range files {
		item := storageComponents.FileItemCompact{
			Filename:   file.Filename,
			MimeType:   pointer.ToString(file.MimeType),
			ObjectPath: file.ObjectPath,
			Size:       pointer.ToInt(file.Size),
			Uid:        file.UID.String(),
		}
		filesList = append(filesList, item)
	}
	response := &storageComponents.FilesListResponse{
		Files: filesList,
	}

	s.responseOK(c, response)
}
