package service

import (
	"context"
	"database/sql"
	"encoding/json"
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
	ItemStore     store.ItemStore
	SearchStore   store.SearchStore
	ItemIndexName string
}

func NewItemService(itemStore store.ItemStore, searchStore store.SearchStore) *ItemService {
	return &ItemService{
		ItemStore:     itemStore,
		SearchStore:   searchStore,
		ItemIndexName: "products_index",
	}
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

	err = s.SearchStore.IndexDocument(ctx, s.ItemIndexName, fmt.Sprintf("%d", createdItem.ID), createdItem)
	if err != nil {
		fmt.Printf("Warning: Failed to index new item %d in Elasticsearch: %v\n", createdItem.ID, err)
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

	err = s.SearchStore.IndexDocument(ctx, s.ItemIndexName, fmt.Sprintf("%d", updatedItem.ID), updatedItem)
	if err != nil {
		fmt.Printf("Warning: Failed to re-index updated item %d in Elasticsearch: %v\n", updatedItem.ID, err)
	}

	return updatedItem, nil
}

func (s *ItemService) DeleteItem(ctx context.Context, id int32) error {
	err := s.ItemStore.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("item not found")
		}
		return fmt.Errorf("failed to delete item from DB: %w", err)
	}

	err = s.SearchStore.DeleteDocument(ctx, s.ItemIndexName, fmt.Sprintf("%d", id))
	if err != nil {
		fmt.Printf("Warning: Failed to delete item %d from Elasticsearch: %v\n", id, err)
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

func (s *ItemService) SearchItems(ctx context.Context, query string, page, pageSize int) (*PaginatedItems, error) {
	searchFields := []string{"name", "description"}

	rawHits, totalCount64, err := s.SearchStore.Search(ctx, s.ItemIndexName, query, searchFields, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to search items in Elasticsearch: %w", err)
	}

	var items []model.Item
	for _, rawHit := range rawHits {
		var item model.Item
		if err := json.Unmarshal(rawHit, &item); err != nil {
			fmt.Printf("Warning: Failed to unmarshal item from Elasticsearch hit: %v\n", err)
			continue
		}
		items = append(items, item)
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
