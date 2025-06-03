package store

import (
	"context"
	"database/sql"
	"fmt"

	"external-backend-go/db/sqlc"
	"external-backend-go/internal/model"
)

type ItemStore interface {
	RepositoryInterface[*model.Item]
}

type itemStore struct {
	*BaseRepository
	queries *sqlc.Queries
}

func NewItemStore(db *sql.DB, queries *sqlc.Queries, baseRepo *BaseRepository) ItemStore {
	return &itemStore{BaseRepository: baseRepo, queries: queries}
}

func (s *itemStore) Create(ctx context.Context, item *model.Item) (*model.Item, error) {
	params := sqlc.CreateItemParams{
		Name:        item.Name,
		Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
	}
	createdItem, err := s.queries.CreateItem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create item in DB: %w", err)
	}
	return &model.Item{
		ID:          createdItem.ID,
		Name:        createdItem.Name,
		Description: createdItem.Description.String,
		CreatedAt:   createdItem.CreatedAt,
		UpdatedAt:   createdItem.UpdatedAt,
	}, nil
}

func (s *itemStore) GetByID(ctx context.Context, id int32) (*model.Item, error) {
	dbItem, err := s.queries.GetItemByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get item by ID from DB: %w", err)
	}
	return &model.Item{
		ID:          dbItem.ID,
		Name:        dbItem.Name,
		Description: dbItem.Description.String,
		CreatedAt:   dbItem.CreatedAt,
		UpdatedAt:   dbItem.UpdatedAt,
	}, nil
}

func (s *itemStore) Update(ctx context.Context, item *model.Item) (*model.Item, error) {
	params := sqlc.UpdateItemParams{
		ID:          item.ID,
		Name:        item.Name,
		Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
	}
	updatedItem, err := s.queries.UpdateItem(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to update item in DB: %w", err)
	}
	return &model.Item{
		ID:          updatedItem.ID,
		Name:        updatedItem.Name,
		Description: updatedItem.Description.String,
		CreatedAt:   updatedItem.CreatedAt,
		UpdatedAt:   updatedItem.UpdatedAt,
	}, nil
}

func (s *itemStore) Delete(ctx context.Context, id int32) error {
	err := s.queries.DeleteItem(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to delete item from DB: %w", err)
	}
	return nil
}

func (s *itemStore) List(ctx context.Context, offset, limit int32) ([]*model.Item, error) {
	params := sqlc.ListItemsParams{
		Offset: offset,
		Limit:  limit,
	}
	dbItems, err := s.queries.ListItems(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get items from DB: %w", err)
	}

	var items []*model.Item
	for _, dbItem := range dbItems {
		items = append(items, &model.Item{
			ID:          dbItem.ID,
			Name:        dbItem.Name,
			Description: dbItem.Description.String,
			CreatedAt:   dbItem.CreatedAt,
			UpdatedAt:   dbItem.UpdatedAt,
		})
	}
	return items, nil
}

func (s *itemStore) Count(ctx context.Context) (int64, error) {
	count, err := s.queries.CountItems(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count items in DB: %w", err)
	}
	return count, nil
}
