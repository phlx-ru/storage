package data

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/phlx-ru/hatchet/logger"
	"github.com/phlx-ru/hatchet/metrics"
	"github.com/phlx-ru/hatchet/watcher"

	"storage/ent"
	"storage/ent/predicate"
)

const (
	metricPrefix = `data.file`
)

type FileRepo struct {
	data    Database
	metric  metrics.Metrics
	logs    *log.Helper
	watcher *watcher.Watcher
}

func NewFileRepo(data Database, logs log.Logger, metric metrics.Metrics) *FileRepo {
	loggerHelper := logger.NewHelper(logs, `ts`, log.DefaultTimestamp, `scope`, metricPrefix)
	return &FileRepo{
		data:    data,
		metric:  metric,
		logs:    loggerHelper,
		watcher: watcher.New(metricPrefix, loggerHelper, metric),
	}
}

func (f *FileRepo) Create(ctx context.Context, file *ent.File) (*ent.File, error) {
	var err error
	defer f.watcher.OnPreparedMethod(`Create`).WithFields(map[string]any{
		"filename":   file.Filename,
		"objectPath": file.ObjectPath,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

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

func (f *FileRepo) Delete(ctx context.Context, uid string) error {
	var err error
	defer f.watcher.OnPreparedMethod(`Delete`).WithFields(map[string]any{
		"uid": uid,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	_, err = f.client(ctx).
		Update().
		Where(fileFilterByUID(uid)).
		SetDeletedAt(time.Now()).
		Save(ctx)

	return err
}

func (f *FileRepo) Restore(ctx context.Context, uid string) error {
	var err error
	defer f.watcher.OnPreparedMethod(`Restore`).WithFields(map[string]any{
		"uid": uid,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	_, err = f.client(ctx).
		Update().
		Where(fileFilterByUID(uid)).
		ClearDeletedAt().
		Save(ctx)

	return err
}

func (f *FileRepo) FindByUID(ctx context.Context, uid string) (*ent.File, error) {
	var err error
	defer f.watcher.OnPreparedMethod(`FindByUID`).WithFields(map[string]any{
		"uid": uid,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	file, err := f.client(ctx).
		Query().
		Where(fileFilterActive()).
		Where(fileFilterByUID(uid)).
		Only(ctx)

	return file, err
}

func (f *FileRepo) FindByUserID(ctx context.Context, userID, limit, offset int) ([]*ent.File, error) {
	var err error
	defer f.watcher.OnPreparedMethod(`FindByUID`).WithFields(map[string]any{
		"userID": userID,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	file, err := f.client(ctx).
		Query().
		Where(fileFilterActive()).
		Where(fileFilterByUserID(userID)).
		Limit(limit).
		Offset(offset).
		All(ctx)

	return file, err
}

func (f *FileRepo) FindByFilename(ctx context.Context, filename string) (*ent.File, error) {
	var err error
	defer f.watcher.OnPreparedMethod(`FindByFilename`).WithFields(map[string]any{
		"filename": filename,
	}).WithIgnoredErrorsChecks([]func(error) bool{
		ent.IsNotFound,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	file, err := f.client(ctx).
		Query().
		Where(fileFilterActive()).
		Where(fileFilterByFilename(filename)).
		First(ctx)

	return file, err
}

func (f *FileRepo) FindByObjectPath(ctx context.Context, objectPath string) (*ent.File, error) {
	var err error
	defer f.watcher.OnPreparedMethod(`FindByObjectPath`).WithFields(map[string]any{
		"objectPath": objectPath,
	}).WithIgnoredErrorsChecks([]func(error) bool{
		ent.IsNotFound,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})

	file, err := f.client(ctx).
		Query().
		Where(fileFilterActive()).
		Where(fileFilterByObjectPath(objectPath)).
		First(ctx)

	return file, err
}

func (f *FileRepo) client(ctx context.Context) *ent.FileClient {
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
