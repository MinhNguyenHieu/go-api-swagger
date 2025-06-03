package store

import (
	"context"
	"external-backend-go/internal/model"
)

type CacheStorage interface {
	Users() UserCache
}

type UserCache interface {
	Get(ctx context.Context, userID string) (*model.User, error)
	Set(ctx context.Context, user *model.User) error
}

type DummyUserCache struct{}

func (d *DummyUserCache) Get(ctx context.Context, userID string) (*model.User, error) {
	return nil, nil
}

func (d *DummyUserCache) Set(ctx context.Context, user *model.User) error {
	return nil
}

type dummyCacheStorage struct {
	userCache UserCache
}

func NewDummyCacheStorage() CacheStorage {
	return &dummyCacheStorage{
		userCache: &DummyUserCache{},
	}
}

func (d *dummyCacheStorage) Users() UserCache {
	return d.userCache
}
