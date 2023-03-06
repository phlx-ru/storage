package minio

import (
	"context"
	"io"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"github.com/phlx-ru/hatchet/logger"
	"github.com/phlx-ru/hatchet/metrics"
	"github.com/phlx-ru/hatchet/watcher"
)

const (
	metricPrefix = `clients.minio`
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

type Minio struct {
	minio      *minio.Client
	bucketName string
	metric     metrics.Metrics
	logger     *log.Helper
	watcher    *watcher.Watcher
}

func New(
	endpoint, bucketLocation, bucketName, accessKeyID, secretAccessKey string,
	metric metrics.Metrics,
	logs log.Logger,
) (*Minio, error) {
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
	loggerHelper := logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, metricPrefix)

	return &Minio{
		minio:      client,
		bucketName: bucketName,
		metric:     metric,
		logger:     loggerHelper,
		watcher:    watcher.New(metricPrefix, loggerHelper, metric),
	}, nil
}

func (c *Minio) Upload(ctx context.Context, filePath string, objectPath string) (minio.UploadInfo, error) {
	var err error
	defer c.watcher.OnPreparedMethod(`Upload`).Results(func() (context.Context, error) {
		return ctx, err
	})

	uploadInfo, err := c.minio.FPutObject(
		ctx,
		c.bucketName,
		objectPath,
		filePath,
		minio.PutObjectOptions{},
	)

	return uploadInfo, err
}

func (c *Minio) Download(ctx context.Context, filePath string, objectPath string) error {
	var err error
	defer c.watcher.OnPreparedMethod(`Download`).Results(func() (context.Context, error) {
		return ctx, err
	})

	err = c.minio.FGetObject(
		ctx,
		c.bucketName,
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

func (c *Minio) Remove(ctx context.Context, objectPath string) error {
	var err error
	defer c.watcher.OnPreparedMethod(`Remove`).Results(func() (context.Context, error) {
		return ctx, err
	})

	err = c.minio.RemoveObject(ctx, c.bucketName, objectPath, minio.RemoveObjectOptions{})

	return err
}

func (c *Minio) UploadFromReader(
	ctx context.Context,
	reader io.Reader,
	size int64,
	contentType string,
	objectPath string,
) (minio.UploadInfo, error) {
	var err error
	defer c.watcher.OnPreparedMethod(`UploadFromReader`).Results(func() (context.Context, error) {
		return ctx, err
	})

	uploadInfo := minio.UploadInfo{}
	err = c.minio.RemoveObject(ctx, c.bucketName, objectPath, minio.RemoveObjectOptions{})
	if err != nil {
		return uploadInfo, err
	}

	err = s3utils.CheckValidObjectName(objectPath)
	if err != nil {
		return uploadInfo, err
	}

	uploadInfo, err = c.minio.PutObject(
		ctx,
		c.bucketName,
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

func (c *Minio) DownloadToWriter(ctx context.Context, writer io.Writer, objectPath string) error {
	var err error
	defer c.watcher.OnPreparedMethod(`DownloadToWriter`).Results(func() (context.Context, error) {
		return ctx, err
	})

	err = s3utils.CheckValidObjectName(objectPath)
	if err != nil {
		return err
	}

	object, err := c.minio.GetObject(ctx, c.bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = object.Close()
	}()

	_, err = io.Copy(writer, object)

	return err
}
