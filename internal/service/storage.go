package service

import (
	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"

	"storage/api/storage"
	storageComponents "storage/api/storage/storage"
	"storage/internal/biz"
	"storage/internal/pkg/logger"
	"storage/internal/pkg/metrics"
	"storage/internal/pkg/request"
	pkgStrings "storage/internal/pkg/strings"
)

const (
	metricPrefix = `service.storage`
	uidMaxLength = 36
)

type StorageService struct {
	usecase *biz.StorageUsecase
	metric  metrics.Metrics
	logger  *log.Helper
}

func NewGatewayService(
	usecase *biz.StorageUsecase,
	metric metrics.Metrics,
	logs log.Logger,
) *StorageService {
	return &StorageService{
		usecase: usecase,
		metric:  metric,
		logger:  logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "service/storage"),
	}
}

func (s *StorageService) GetSwagger(c *gin.Context) {
	swagger, err := storage.GetSwagger()
	if err != nil {
		s.responseError(c, errors.InternalServer(`internal_error`, err.Error()))
		return
	}
	swagger.InternalizeRefs(c.Request.Context(), nil) // TODO Bug with nested refs and properties includes
	s.responseOK(c, swagger)
}

func (s *StorageService) Download(c *gin.Context, uid storageComponents.Uid) {
	method := `download`
	defer s.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { s.postProcess(c.Request.Context(), method, err) }()

	if uid == "" {
		s.responseError(c, errors.BadRequest(`empty_uid`, `UID is empty`))
		return
	}

	if len(uid) > uidMaxLength {
		s.responseError(c, errors.BadRequest(`ruined_uid`, `UID must have 36 symbols or less`))
		return
	}

	err = s.usecase.Download(c.Request.Context(), uid, c.Writer)
	if err != nil {
		s.responseError(c, err)
	}
}

func (s *StorageService) Upload(c *gin.Context, params storage.UploadParams) {
	method := `upload`
	defer s.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { s.postProcess(c.Request.Context(), method, err) }()

	filename := c.Request.URL.Query().Get(`filename`)
	if filename == "" {
		s.responseValidationError(c, errors.Unauthorized(`filename_unspecified`, `filename is not specified`))
		return
	}

	token := request.AuthTokenFromContext(c.Request.Context())
	if params.AuthToken != nil {
		token = *params.AuthToken
	}
	if token == "" {
		s.responseError(c, errors.Unauthorized(`auth_token_unspecified`, `not found auth token in request`))
		return
	}

	uploadFile := &biz.UploadFile{
		Reader:   c.Request.Body,
		Size:     c.Request.ContentLength,
		Filename: filename,
	}

	file, err := s.usecase.Upload(c.Request.Context(), uploadFile, token)
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

func (s *StorageService) FilesList(c *gin.Context, params storage.FilesListParams) {
	method := `filesList`
	defer s.metric.NewTiming().Send(pkgStrings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { s.postProcess(c.Request.Context(), method, err) }()

	token := request.AuthTokenFromContext(c.Request.Context())
	if params.AuthToken != nil {
		token = *params.AuthToken
	}
	if token == "" {
		s.responseError(c, errors.Unauthorized(`auth_token_unspecified`, `not found auth token in request`))
		return
	}

	files, err := s.usecase.FilesList(c.Request.Context(), token)
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
