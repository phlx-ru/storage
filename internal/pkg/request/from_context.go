package request

import (
	"context"
	"strings"

	"github.com/go-kratos/gin"
)

const (
	cookieAuthTokenKey = `auth_token`
)

func AuthTokenFromContext(ctx context.Context) string {
	c, ok := gin.FromGinContext(ctx)
	if !ok {
		return ""
	}
	a := c.GetHeader(`Authorization`)
	if a == "" {
		return ""
	}
	authHeader := c.GetHeader(`Authorization`)
	if authHeader != "" && strings.HasPrefix(authHeader, `Bearer `) {
		return strings.ReplaceAll(authHeader, `Bearer `, ``)
	}
	token, _ := c.Cookie(cookieAuthTokenKey)
	return token
}
