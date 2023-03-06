package biz

import (
	"context"
	"errors"

	"github.com/phlx-ru/hatchet/gin"

	v1 "storage/api/storage/v1"
	"storage/internal/clients/auth"
	"storage/internal/pkg/texts"
)

const (
	userTypeAdmin = `admin`
)

type AuthenticatedUser struct {
	auth *auth.CheckResult
}

//go:generate moq -out authentication_moq_test.go . AuthChecker

type AuthChecker interface {
	Check(ctx context.Context, token string) (*auth.CheckResult, error)
}

func (s *StorageUsecase) isIntegrations(ctx context.Context) bool {
	return IsIntegrations(ctx, s.auth.Jwt.Secret)
}

func IsIntegrations(ctx context.Context, secret string) bool {
	return gin.CheckIntegrationsTokenFromContext(ctx, secret)
}

func (s *StorageUsecase) user(ctx context.Context) (*AuthenticatedUser, error) {
	return User(ctx, s.authClient)
}

func User(ctx context.Context, authClient AuthChecker) (*AuthenticatedUser, error) {
	token := gin.AuthTokenFromContext(ctx)
	if token == "" {
		return nil, v1.ErrorUnauthorized(texts.AccessDenied)
	}
	check, err := authClient.Check(ctx, token)
	if err != nil {
		if errors.Is(err, auth.ErrSessionExpiredOrNotFound) {
			return nil, v1.ErrorUnauthorized(texts.AuthSessionExpired)
		}
		return nil, v1.ErrorInternalError(err.Error())
	}
	if check.User == nil || check.Session == nil {
		return nil, v1.ErrorUnauthorized(texts.AuthSessionUnknown)
	}
	return &AuthenticatedUser{
		auth: check,
	}, nil
}

func (a *AuthenticatedUser) ID() int64 {
	if a.auth == nil || a.auth.User == nil {
		return 0
	}
	return a.auth.User.ID
}

func (a *AuthenticatedUser) AuthUser() *auth.User {
	if a.auth == nil {
		return nil
	}
	return a.auth.User
}

func (a *AuthenticatedUser) AuthSession() *auth.Session {
	if a.auth == nil {
		return nil
	}
	return a.auth.Session
}

func (a *AuthenticatedUser) IsAdmin() bool {
	if a.auth == nil || a.auth.User == nil {
		return false
	}
	return a.auth.User.Type == userTypeAdmin
}
