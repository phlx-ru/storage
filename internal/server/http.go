package server

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kratosHTTP "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/phlx-ru/hatchet/metrics"
	"github.com/phlx-ru/hatchet/middlewares"

	"storage/internal/conf"
	"storage/internal/service"
	storage "storage/schema"
)

const (
	metricPrefix = `server`
)

// NewHTTPServer new HTTP server.
func NewHTTPServer(
	c *conf.Server,
	ss *service.StorageService,
	metric metrics.Metrics,
) *kratosHTTP.Server {
	var opts = []kratosHTTP.ServerOption{
		kratosHTTP.Timeout(c.Http.Timeout.AsDuration()),
	}
	if c.Http.Network != "" {
		opts = append(opts, kratosHTTP.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, kratosHTTP.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, kratosHTTP.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := kratosHTTP.NewServer(opts...)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(
		gin.LoggerWithConfig(
			gin.LoggerConfig{
				Formatter: LogFormatter,
				Output:    gin.DefaultWriter,
				SkipPaths: []string{
					`/swagger/`,
					`/swagger/swagger-ui.css`,
					`/swagger/swagger-ui.css.map`,
					`/swagger/swagger-ui-bundle.js`,
					`/swagger/swagger-ui-bundle.js.map`,
					`/swagger/swagger.yaml`,
					`/swagger/favicon.ico`,
					`/form/`,
					`/form/favicon.ico`,
				},
			},
		),
		gin.Recovery(),
		kgin.Middlewares(
			middlewares.Duration(metric, metricPrefix),
			tracing.Server(),
			recovery.Recovery(),
		),
		cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders: []string{
				"Authorization",
				"Origin",
				"Content-Length",
				"Content-Type",
				"Content-Disposition",
				"Accept",
				"Accept-Language",
				"Accept-Range",
				"Range",
			},
		}),
	)

	router.Static(`/form`, `./static/form`)
	router.Static(`/swagger`, `./static/swagger`)

	router.GET(`/api/swagger`, ss.GetSwagger)

	storage.RegisterHandlersWithOptions(
		router,
		ss,
		storage.GinServerOptions{
			BaseURL:     ``, // c.BaseUrl
			Middlewares: []storage.MiddlewareFunc{},
		},
	)

	srv.HandlePrefix(`/`, router)

	return srv
}

func LogFormatter(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency -= param.Latency % time.Second
	}
	return fmt.Sprintf(
		"ACCESS ts=%v status=%d latency=%v client.ip=%s method=%s path=%-7s error=%#v\n",
		param.TimeStamp.Format(time.RFC3339),
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}
