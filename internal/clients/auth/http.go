package auth

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/phlx-ru/hatchet/logger"
	"github.com/phlx-ru/hatchet/metrics"
	"github.com/phlx-ru/hatchet/watcher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "storage/api/auth/v1"
)

const (
	metricPrefix = `clients.auth`
)

var (
	ErrSessionExpiredOrNotFound = errors.New(`session expired or not found`)
)

type Client interface {
	Check(ctx context.Context, token string) (*CheckResult, error)
}

type Auth struct {
	client  v1.AuthClient
	metric  metrics.Metrics
	logger  *log.Helper
	watcher *watcher.Watcher
}

func New(
	ctx context.Context,
	endpoint string,
	timeout time.Duration,
	metric metrics.Metrics,
	logs log.Logger,
) (*Auth, error) {
	client, err := Default(ctx, endpoint, timeout)
	if err != nil {
		return nil, err
	}
	loggerHelper := logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, metricPrefix)
	return &Auth{
		client:  v1.NewAuthClient(client),
		metric:  metric,
		logger:  loggerHelper,
		watcher: watcher.New(metricPrefix, loggerHelper, metric),
	}, nil
}

type User struct {
	ID          int64   `json:"id"`
	Type        string  `json:"type"`
	DisplayName string  `json:"displayName"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
}

type Session struct {
	Until     time.Time `json:"until"`
	IP        *string   `json:"IP,omitempty"`
	UserAgent *string   `json:"userAgent,omitempty"`
	DeviceID  *string   `json:"deviceId,omitempty"`
}

type CheckResult struct {
	User    *User    `json:"user"`
	Session *Session `json:"session"`
}

func (a *Auth) Check(ctx context.Context, token string) (*CheckResult, error) {
	var err error
	defer a.watcher.OnPreparedMethod(`Check`).Results(func() (context.Context, error) {
		return ctx, err
	})
	res, err := a.client.Check(ctx, &v1.CheckRequest{Token: token})
	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			if statusErr.Code() == codes.NotFound {
				return nil, ErrSessionExpiredOrNotFound
			}
		}
		return nil, err
	}
	return &CheckResult{
		User: &User{
			ID:          res.User.Id,
			Type:        res.User.Type,
			DisplayName: res.User.DisplayName,
			Email:       res.User.Email,
			Phone:       res.User.Phone,
		},
		Session: &Session{
			Until:     res.Session.Until.AsTime(),
			IP:        res.Session.Ip,
			UserAgent: res.Session.UserAgent,
			DeviceID:  res.Session.DeviceId,
		},
	}, nil
}
