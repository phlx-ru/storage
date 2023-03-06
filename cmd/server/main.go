package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"
	pkgConfig "github.com/phlx-ru/hatchet/config"
	"github.com/phlx-ru/hatchet/logger"
	"github.com/phlx-ru/hatchet/metrics"
	"github.com/phlx-ru/hatchet/runtime"

	"storage/internal/clients/auth"
	"storage/internal/clients/minio"
	"storage/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = `storage-server`
	// Version is the version of the compiled software.
	Version = `1.1.1`
	// flagconf is the config flag.
	flagconf string
	// dotenv is loaded from config path .env file
	dotenv string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&dotenv, "dotenv", ".env", ".env file, eg: -dotenv .env")
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()

	var err error

	ctx := context.Background()

	envPath := path.Join(flagconf, dotenv)
	err = godotenv.Overload(envPath)
	if err != nil {
		return err
	}

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(pkgConfig.EnvReplaceDecoder),
	)
	defer func() {
		_ = c.Close()
	}()

	if err = c.Load(); err != nil {
		return err
	}

	var bc conf.Bootstrap
	if err = c.Scan(&bc); err != nil {
		return err
	}

	logs := logger.New(id, Name, Version, bc.Log.Level, bc.Env, bc.Sentry.Level)

	metric, err := metrics.New(bc.Metrics.Address, Name, bc.Metrics.Mute)
	if err != nil {
		return err
	}
	defer metric.Close()
	metric.Increment("starts.count")

	database, cleanup, err := wireData(bc.Data, logs)
	if err != nil {
		return err
	}
	defer cleanup()

	go database.CollectDatabaseMetrics(ctx, metric)
	go runtime.CollectGoMetrics(ctx, metric)

	if err = database.Prepare(ctx, bc.Data.Database.Migrate); err != nil {
		return err
	}

	authConf := bc.Client.Grpc.Auth
	authClient, err := auth.New(ctx, authConf.Endpoint, authConf.Timeout.AsDuration(), metric, logs)
	if err != nil {
		return err
	}

	s3 := selectS3Config(bc.S3)
	minioClient, err := minio.New(s3.Endpoint, s3.BucketLocation, s3.BucketName, s3.AccessKeyID, s3.SecretAccessKey, metric, logs)
	if err != nil {
		return err
	}

	app, err := wireApp(ctx, database, bc.Server, bc.Auth, authClient, minioClient, metric, logs)
	if err != nil {
		panic(err)
	}

	return app.Run()
}

type S3Config struct {
	Endpoint        string
	BucketLocation  string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
}

func selectS3Config(s3 *conf.S3) *S3Config {
	switch s3.Current {
	case `vk`:
		return &S3Config{
			Endpoint:        s3.Vk.Endpoint,
			BucketLocation:  s3.Vk.BucketLocation,
			BucketName:      s3.Vk.BucketName,
			AccessKeyID:     s3.Vk.AccessKeyID,
			SecretAccessKey: s3.Vk.SecretAccessKey,
		}
	case `yandex`:
		return &S3Config{
			Endpoint:        s3.Yandex.Endpoint,
			BucketLocation:  s3.Yandex.BucketLocation,
			BucketName:      s3.Yandex.BucketName,
			AccessKeyID:     s3.Yandex.AccessKeyID,
			SecretAccessKey: s3.Yandex.SecretAccessKey,
		}
	}
	panic(fmt.Sprintf(`unknown s3 current value: %s`, s3.Current))
}

func newApp(ctx context.Context, logger log.Logger, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Context(ctx),
		kratos.Server(
			hs,
		),
	)
}
