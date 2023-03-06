package biz

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	kratosGin "github.com/go-kratos/gin"
	"github.com/stretchr/testify/require"

	v1 "storage/api/storage/v1"
	"storage/internal/clients/auth"
)

var (
	checkResultAdmin = &auth.CheckResult{
		User: &auth.User{
			ID:          1,
			Type:        "admin",
			DisplayName: "I am the Admin",
			Email:       pointer.ToString("theadmin@nonexistent.test"),
		},
		Session: &auth.Session{
			Until:     time.Now().Add(48 * time.Hour),
			IP:        pointer.ToString("10.0.0.10"),
			UserAgent: pointer.ToString("Golang Test User-Agent"),
			DeviceID:  nil,
		},
	}

	checkResultDispatcher = &auth.CheckResult{
		User: &auth.User{
			ID:          2,
			Type:        "dispatcher",
			DisplayName: "I am only Dispatcher",
			Email:       pointer.ToString("thedispatcher@nonexistent.test"),
		},
		Session: &auth.Session{
			Until:     time.Now().Add(24 * time.Hour),
			IP:        pointer.ToString("10.0.0.10"),
			UserAgent: pointer.ToString("Golang Test User-Agent"),
			DeviceID:  nil,
		},
	}

	checkResultMalformed = &auth.CheckResult{
		User: &auth.User{
			ID:          3,
			Type:        "driver",
			DisplayName: "I am malformed",
			Email:       pointer.ToString("malformed@nonexistent.test"),
		},
		Session: nil,
	}
)

func TestUser(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name           string
		ctx            context.Context
		authClientMock func() *AuthCheckerMock
		expectedUser   *AuthenticatedUser
		expectedError  func(error) bool
	}{
		{
			name: "admin",
			ctx: kratosGin.NewGinContext(ctx, &gin.Context{
				Request: &http.Request{
					Header: map[string][]string{
						"Authorization": {"Bearer let0me0in"},
					},
				},
			}),
			authClientMock: func() *AuthCheckerMock {
				return &AuthCheckerMock{
					CheckFunc: func(ctx context.Context, token string) (*auth.CheckResult, error) {
						if token != "let0me0in" {
							return nil, auth.ErrSessionExpiredOrNotFound
						}
						return checkResultAdmin, nil
					},
				}
			},
			expectedUser: &AuthenticatedUser{
				auth: checkResultAdmin,
			},
		},
		{
			name: "dispatcher",
			ctx: kratosGin.NewGinContext(ctx, &gin.Context{
				Request: &http.Request{
					Header: map[string][]string{
						"Authorization": {"Bearer i0am0dispatcher"},
					},
				},
			}),
			authClientMock: func() *AuthCheckerMock {
				return &AuthCheckerMock{
					CheckFunc: func(ctx context.Context, token string) (*auth.CheckResult, error) {
						if token != "i0am0dispatcher" {
							return nil, auth.ErrSessionExpiredOrNotFound
						}
						return checkResultDispatcher, nil
					},
				}
			},
			expectedUser: &AuthenticatedUser{
				auth: checkResultDispatcher,
			},
		},
		{
			name: "unauthorized",
			ctx: kratosGin.NewGinContext(ctx, &gin.Context{
				Request: &http.Request{
					Header: map[string][]string{
						"Authorization": {"Bearer nonexistent"},
					},
				},
			}),
			authClientMock: func() *AuthCheckerMock {
				return &AuthCheckerMock{
					CheckFunc: func(ctx context.Context, token string) (*auth.CheckResult, error) {
						return nil, auth.ErrSessionExpiredOrNotFound
					},
				}
			},
			expectedError: v1.IsUnauthorized,
		},
		{
			name: "internal_error",
			ctx: kratosGin.NewGinContext(ctx, &gin.Context{
				Request: &http.Request{
					Header: map[string][]string{
						"Authorization": {"Bearer nonexistent"},
					},
				},
			}),
			authClientMock: func() *AuthCheckerMock {
				return &AuthCheckerMock{
					CheckFunc: func(ctx context.Context, token string) (*auth.CheckResult, error) {
						return nil, errors.New("some internal error")
					},
				}
			},
			expectedError: v1.IsInternalError,
		},
		{
			name: "empty_session",
			ctx: kratosGin.NewGinContext(ctx, &gin.Context{
				Request: &http.Request{
					Header: map[string][]string{
						"Authorization": {"Bearer dont0matter"},
					},
				},
			}),
			authClientMock: func() *AuthCheckerMock {
				return &AuthCheckerMock{
					CheckFunc: func(ctx context.Context, token string) (*auth.CheckResult, error) {
						return checkResultMalformed, nil
					},
				}
			},
			expectedError: v1.IsUnauthorized,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualUser, err := User(testCase.ctx, testCase.authClientMock())
			if testCase.expectedError != nil {
				require.Error(t, err)
				require.True(t, testCase.expectedError(err))
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedUser, actualUser)
			}
		})
	}
}
