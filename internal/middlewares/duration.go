package middlewares

import (
	"context"
	"strings"

	"storage/internal/pkg/logger"
	"storage/internal/pkg/metrics"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func Duration(metric metrics.Metrics, logs log.Logger) middleware.Middleware {
	lg := logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "middlewares/duration")
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			slug := metricSlug(ctx)
			timing := metric.NewTiming()
			defer func() {
				if slug == "" {
					lg.Warn("failed to parse metric slug from context")
				} else {
					lg.Debugf("slug for duration metric: [%s]", slug)
					lg.Debugf("duration metric timing: %d ms", timing.Duration().Milliseconds())
					timing.Send(slug)
				}
			}()
			return handler(ctx, req)
		}
	}
}

func metricSlug(ctx context.Context) string {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return ""
	}
	h, ok := tr.(*http.Transport)
	if !ok {
		return ""
	}

	slug := strings.ToLower(h.Request().URL.Path)
	slug = strings.ReplaceAll(strings.ToLower(slug), ".", "_")
	slug = strings.ReplaceAll(slug, "/", "_")
	kind := strings.ToLower(h.Request().Method)
	return kind + slug
}
