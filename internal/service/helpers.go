package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/phlx-ru/hatchet/validate"

	v1 "storage/api/storage/v1"
	"storage/internal/pkg/texts"
)

func (s *StorageService) validate(something any) error {
	validationErrors, err := validate.Struct(something)
	if err != nil {
		return errors.InternalServer(`validator_fails`, err.Error())
	}
	if validationErrors != nil {
		metadata := validate.AsCustomValidationTranslations(validationErrors)
		return v1.ErrorValidationFailed(texts.ValidationFailed).WithMetadata(metadata)
	}
	return nil
}

func (s *StorageService) responseError(c *gin.Context, err error) {
	c.Header(`Content-Type`, `application/json`) // Bugfix for error render
	kgin.Error(c, err)
}

func (s *StorageService) responseValidationError(c *gin.Context, err error) {
	validationErr := err
	if !v1.IsValidationFailed(err) {
		validationErr = v1.ErrorValidationFailed(err.Error())
	}
	s.responseError(c, validationErr)
}

func (s *StorageService) responseOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}
