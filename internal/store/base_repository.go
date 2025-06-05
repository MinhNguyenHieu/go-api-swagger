package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"external-backend-go/internal/logger"
)

type GenericEntity interface {
	GetID() int32
	SetID(id int32)
	GetCreatedAt() time.Time
	SetCreatedAt(t time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(t time.Time)
}

type RepositoryInterface[T GenericEntity] interface {
	Create(ctx context.Context, entity T) (T, error)
	GetByID(ctx context.Context, id int32) (T, error)
	Update(ctx context.Context, entity T) (T, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, offset, limit int32) ([]T, error)
	Count(ctx context.Context) (int64, error)
}

type BaseRepository struct {
	DB     *sql.DB
	Logger *logger.Logger
}

func NewBaseRepository(db *sql.DB, appLogger *logger.Logger) *BaseRepository {
	return &BaseRepository{
		DB:     db,
		Logger: appLogger,
	}
}

func (b *BaseRepository) Count(ctx context.Context) (int64, error) {
	b.Logger.Warn("Generic BaseRepository.Count called. This should typically be overridden by concrete stores.")
	return 0, fmt.Errorf("generic count not implemented for BaseRepository; override in concrete store")
}

func (b *BaseRepository) List(ctx context.Context, offset, limit int32) ([]GenericEntity, error) {
	b.Logger.Warn("Generic BaseRepository.List called. This should typically be overridden by concrete stores.")
	return nil, fmt.Errorf("generic list not implemented for BaseRepository; override in concrete store")
}

func (b *BaseRepository) Create(ctx context.Context, entity GenericEntity) (GenericEntity, error) {
	b.Logger.Warn("Generic BaseRepository.Create called. This should typically be overridden by concrete stores.")
	return nil, fmt.Errorf("generic create not implemented for BaseRepository; override in concrete store")
}

func (b *BaseRepository) GetByID(ctx context.Context, id int32) (GenericEntity, error) {
	b.Logger.Warn("Generic BaseRepository.GetByID called. This should typically be overridden by concrete stores.")
	return nil, fmt.Errorf("generic GetByID not implemented for BaseRepository; override in concrete store")
}

func (b *BaseRepository) Update(ctx context.Context, entity GenericEntity) (GenericEntity, error) {
	b.Logger.Warn("Generic BaseRepository.Update called. This should typically be overridden by concrete stores.")
	return nil, fmt.Errorf("generic Update not implemented for BaseRepository; override in concrete store")
}

func (b *BaseRepository) Delete(ctx context.Context, id int32) error {
	b.Logger.Warn("Generic BaseRepository.Delete called. This should typically be overridden by concrete stores.")
	return fmt.Errorf("generic Delete not implemented for BaseRepository; override in concrete store")
}
