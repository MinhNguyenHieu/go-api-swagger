package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"external-backend-go/internal/model"
	"external-backend-go/internal/store"
)

type PaginatedItems struct {
	Items      []model.Item `json:"items"`
	TotalCount int          `json:"totalCount"`
	Page       int          `json:"page"`
	PageSize   int          `json:"pageSize"`
	TotalPages int          `json:"totalPages"`
}

type ItemService struct {
	ItemStore store.ItemStore
}

func NewItemService(itemStore store.ItemStore) *ItemService {
	return &ItemService{ItemStore: itemStore}
}

func (s *ItemService) CreateItem(ctx context.Context, name, description string) (*model.Item, error) {
	item := &model.Item{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdItem, err := s.ItemStore.Create(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}
	return createdItem, nil
}

func (s *ItemService) GetItemByID(ctx context.Context, id int32) (*model.Item, error) {
	item, err := s.ItemStore.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("item not found")
		}
		return nil, fmt.Errorf("failed to get item by ID: %w", err)
	}
	return item, nil
}

func (s *ItemService) UpdateItem(ctx context.Context, id int32, name, description string) (*model.Item, error) {
	existingItem, err := s.ItemStore.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("item not found")
		}
		return nil, fmt.Errorf("failed to retrieve item for update: %w", err)
	}

	existingItem.Name = name
	existingItem.Description = description
	existingItem.UpdatedAt = time.Now()

	updatedItem, err := s.ItemStore.Update(ctx, existingItem)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}
	return updatedItem, nil
}

func (s *ItemService) DeleteItem(ctx context.Context, id int32) error {
	err := s.ItemStore.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("item not found")
		}
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

func (s *ItemService) GetItems(ctx context.Context, page, pageSize int) (*PaginatedItems, error) {
	offset := (page - 1) * pageSize
	ptrItems, err := s.ItemStore.List(ctx, int32(offset), int32(pageSize))
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	var items []model.Item
	for _, itemPtr := range ptrItems {
		items = append(items, *itemPtr)
	}

	totalCount64, err := s.ItemStore.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count items: %w", err)
	}
	totalCount := int(totalCount64)

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 && totalCount > 0 {
		totalPages = 1
	}

	return &PaginatedItems{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
