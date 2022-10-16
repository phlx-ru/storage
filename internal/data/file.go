package data

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"

	"storage/ent"
	"storage/ent/predicate"
	"storage/internal/biz"
	"storage/internal/pkg/logger"
	"storage/internal/pkg/metrics"
	"storage/internal/pkg/strings"
)

const (
	metricPrefix = `data.file`
)

type fileRepo struct {
	data   Database
	metric metrics.Metrics
	logs   *log.Helper
}

func NewFileRepo(data Database, logs log.Logger, metric metrics.Metrics) biz.FileRepo {
	return &fileRepo{
		data:   data,
		metric: metric,
		logs:   logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, `data/file`),
	}
}

func (f *fileRepo) postProcess(ctx context.Context, method string, err error) {
	if err != nil {
		f.logs.WithContext(ctx).Errorf(`file repo method %s failed: %v`, method, err)
		f.metric.Increment(strings.Metric(metricPrefix, method, `failure`))
	} else {
		f.metric.Increment(strings.Metric(metricPrefix, method, `success`))
	}
}

func (f *fileRepo) Create(ctx context.Context, file *ent.File) (*ent.File, error) {
	method := `create`
	defer f.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { f.postProcess(ctx, method, err) }()

	creating := f.client(ctx).Create().
		SetUserID(file.UserID).
		SetFilename(file.Filename).
		SetObjectPath(file.ObjectPath).
		SetSize(file.Size).
		SetMimeType(file.MimeType)

	if file.DeletedAt != nil {
		creating.SetDeletedAt(*file.DeletedAt)
	}

	created, err := creating.Save(ctx)

	return created, err
}

func (f *fileRepo) Delete(ctx context.Context, uid string) error {
	method := `delete`
	defer f.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { f.postProcess(ctx, method, err) }()

	_, err = f.client(ctx).
		Update().
		Where(fileFilterByUID(uid)).
		SetDeletedAt(time.Now()).
		Save(ctx)

	return err
}

func (f *fileRepo) Restore(ctx context.Context, uid string) error {
	method := `restore`
	defer f.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { f.postProcess(ctx, method, err) }()

	_, err = f.client(ctx).
		Update().
		Where(fileFilterByUID(uid)).
		ClearDeletedAt().
		Save(ctx)

	return err
}

func (f *fileRepo) FindByUID(ctx context.Context, uid string) (*ent.File, error) {
	method := `findByUID`
	defer f.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { f.postProcess(ctx, method, err) }()

	file, err := f.client(ctx).
		Query().
		Where(fileFilterActive()).
		Where(fileFilterByUID(uid)).
		Only(ctx)

	return file, err
}

func (f *fileRepo) FindByUserID(ctx context.Context, userID, limit, offset int) ([]*ent.File, error) {
	method := `findByUserID`
	defer f.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { f.postProcess(ctx, method, err) }()

	file, err := f.client(ctx).
		Query().
		Where(fileFilterActive()).
		Where(fileFilterByUserID(userID)).
		Limit(limit).
		Offset(offset).
		All(ctx)

	return file, err
}

func (f *fileRepo) FindByFilename(ctx context.Context, filename string) ([]*ent.File, error) {
	method := `findByFilename`
	defer f.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { f.postProcess(ctx, method, err) }()

	file, err := f.client(ctx).
		Query().
		Where(fileFilterActive()).
		Where(fileFilterByFilename(filename)).
		All(ctx)

	return file, err
}

func (f *fileRepo) FindByObjectPath(ctx context.Context, objectPath string) ([]*ent.File, error) {
	method := `findByObjectPath`
	defer f.metric.NewTiming().Send(strings.Metric(metricPrefix, method, `timings`))
	var err error
	defer func() { f.postProcess(ctx, method, err) }()

	file, err := f.client(ctx).
		Query().
		Where(fileFilterActive()).
		Where(fileFilterByObjectPath(objectPath)).
		All(ctx)

	return file, err
}

func (f *fileRepo) client(ctx context.Context) *ent.FileClient {
	return client(f.data)(ctx).File
}

func fileFilterActive() predicate.File {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().IsNull(`deleted_at`))
	}
}

func fileFilterByUID(uid string) predicate.File {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().EQ(`uid`, uid))
	}
}

func fileFilterByUserID(userID int) predicate.File {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().EQ(`user_id`, userID))
	}
}

func fileFilterByFilename(filename string) predicate.File {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().EQ(`filename`, filename))
	}
}

func fileFilterByObjectPath(objectPath string) predicate.File {
	return func(selector *sql.Selector) {
		selector.Where(sql.P().EQ(`object_path`, objectPath))
	}
}
