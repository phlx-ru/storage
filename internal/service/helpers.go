package service

import (
	"context"
	"net/http"

	pkgStrings "storage/internal/pkg/strings"
	"storage/internal/pkg/validate"

	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

func (s *StorageService) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		s.logger.WithContext(ctx).Errorf(`method %s failed: %v`, method, err)
		s.metric.Increment(pkgStrings.Metric(method, `failure`))
	} else {
		s.metric.Increment(pkgStrings.Metric(method, `success`))
	}
}

func (s *StorageService) validate(a any) error {
	validationErrors, err := validate.Struct(a)
	if err != nil {
		return errors.InternalServer(`validator_fails`, err.Error())
	}
	if validationErrors != nil {
		metadata := validate.AsCustomValidationTranslations(validationErrors)
		return errors.BadRequest(`validation_error`, `incorrect input`).WithMetadata(metadata)
	}
	return nil
}

func (s *StorageService) responseError(c *gin.Context, err error) {
	c.Header(`Content-Type`, `application/json`) // Bugfix for error render
	kgin.Error(c, err)
}

func (s *StorageService) responseValidationError(c *gin.Context, err error) {
	validationErr := errors.BadRequest(`validation_error`, err.Error())
	s.responseError(c, validationErr)
}

func (s *StorageService) responseOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func (s *StorageService) responseNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, ``)
}
