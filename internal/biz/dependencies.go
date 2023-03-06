package biz

import (
	"context"

	"github.com/google/wire"

	"storage/ent"
	"storage/internal/data"
)

var (
	BindFileRepository = wire.Bind(new(fileRepository), new(*data.FileRepo))
)

//go:generate moq -out dependencies_moq_test.go . fileRepository

type fileRepository interface {
	Create(ctx context.Context, file *ent.File) (*ent.File, error)
	Delete(ctx context.Context, uid string) error
	Restore(ctx context.Context, uid string) error
	FindByUID(ctx context.Context, uid string) (*ent.File, error)
	FindByUserID(ctx context.Context, userID, limit, offset int) ([]*ent.File, error)
	FindByFilename(ctx context.Context, filename string) (*ent.File, error)
	FindByObjectPath(ctx context.Context, objectPath string) (*ent.File, error)
}
