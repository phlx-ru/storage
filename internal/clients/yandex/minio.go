package yandex

import (
	"context"
	"io"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"

	"storage/internal/pkg/logger"
	"storage/internal/pkg/metrics"
	"storage/internal/pkg/strings"
)

const (
	endpoint       = `cargo-staging-storage.storage.yandexcloud.net`
	bucketLocation = `ru-central1`
	bucketName     = `main`

	metricPrefix = `clients.yandex`
)

type Client interface {
	Upload(ctx context.Context, filePath string, objectPath string) (minio.UploadInfo, error)
	Download(ctx context.Context, filePath string, objectPath string) error
	Remove(ctx context.Context, objectPath string) error
	UploadFromReader(
		ctx context.Context,
		reader io.Reader,
		size int64,
		contentType string,
		objectPath string,
	) (minio.UploadInfo, error)
	DownloadToWriter(ctx context.Context, writer io.Writer, objectPath string) error
}

type Yandex struct {
	minio  *minio.Client
	metric metrics.Metrics
	logger *log.Helper
}

func New(
	accessKeyID, secretAccessKey string,
	metric metrics.Metrics,
	logs log.Logger,
) (*Yandex, error) {
	client, err := minio.New(
		endpoint,
		&minio.Options{
			Creds:        credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure:       false,
			Region:       bucketLocation,
			BucketLookup: minio.BucketLookupAuto,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Yandex{
		minio:  client,
		metric: metric,
		logger: logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "clients/yandex"),
	}, nil
}

func (c *Yandex) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		c.logger.WithContext(ctx).Errorf(`client yandex method %s failed: %v`, method, err)
		c.metric.Increment(strings.Metric(metricPrefix, method, `failure`))
	} else {
		c.metric.Increment(strings.Metric(metricPrefix, method, `success`))
	}
}

func (c *Yandex) Upload(ctx context.Context, filePath string, objectPath string) (minio.UploadInfo, error) {
	method := `upload`
	defer c.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { c.postProcess(ctx, method, err) }()

	uploadInfo, err := c.minio.FPutObject(
		ctx,
		bucketName,
		objectPath,
		filePath,
		minio.PutObjectOptions{},
	)

	return uploadInfo, err
}

func (c *Yandex) Download(ctx context.Context, filePath string, objectPath string) error {
	method := `download`
	defer c.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { c.postProcess(ctx, method, err) }()

	err = c.minio.FGetObject(
		ctx,
		bucketName,
		objectPath,
		filePath,
		minio.GetObjectOptions{
			ServerSideEncryption: nil,
			VersionID:            "",
			PartNumber:           0,
			Checksum:             false,
			Internal:             minio.AdvancedGetOptions{},
		},
	)

	return err
}

func (c *Yandex) Remove(ctx context.Context, objectPath string) error {
	method := `remove`
	defer c.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { c.postProcess(ctx, method, err) }()

	err = c.minio.RemoveObject(ctx, bucketName, objectPath, minio.RemoveObjectOptions{})

	return err
}

func (c *Yandex) UploadFromReader(
	ctx context.Context,
	reader io.Reader,
	size int64,
	contentType string,
	objectPath string,
) (minio.UploadInfo, error) {
	method := `uploadFromReader`
	defer c.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { c.postProcess(ctx, method, err) }()

	uploadInfo := minio.UploadInfo{}
	err = c.minio.RemoveObject(ctx, bucketName, objectPath, minio.RemoveObjectOptions{})
	if err != nil {
		return uploadInfo, err
	}

	err = s3utils.CheckValidObjectName(objectPath)
	if err != nil {
		return uploadInfo, err
	}

	uploadInfo, err = c.minio.PutObject(
		ctx,
		bucketName,
		objectPath,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType:        contentType,
			ContentEncoding:    "", // TODO
			ContentDisposition: "", // TODO
			ContentLanguage:    "", // TODO
		},
	)

	return uploadInfo, err
}

func (c *Yandex) DownloadToWriter(ctx context.Context, writer io.Writer, objectPath string) error {
	method := `downloadToWriter`
	defer c.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { c.postProcess(ctx, method, err) }()

	err = s3utils.CheckValidObjectName(objectPath)
	if err != nil {
		return err
	}

	object, err := c.minio.GetObject(ctx, bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = object.Close()
	}()

	_, err = io.Copy(writer, object)

	return err
}
