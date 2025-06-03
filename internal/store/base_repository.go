package store

import (
	"context"
	"database/sql"
	"time"

	"external-backend-go/internal/logger"
)

type ModelWithID interface {
	GetID() int32
	SetID(id int32)
	GetCreatedAt() time.Time
	SetCreatedAt(t time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(t time.Time)
}

type RepositoryInterface[T ModelWithID] interface {
	Create(ctx context.Context, entity T) (T, error)
	GetByID(ctx context.Context, id int32) (T, error)
	Update(ctx context.Context, entity T) (T, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, offset, limit int32) ([]T, error)
	Count(ctx context.Context) (int64, error)
}

type BaseRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewBaseRepository(db *sql.DB, logger *logger.Logger) *BaseRepository {
	return &BaseRepository{db: db, logger: logger}
}
