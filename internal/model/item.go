package model

import (
	"time"
)

type Item struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (i *Item) GetID() int32 {
	return i.ID
}

func (i *Item) SetID(id int32) {
	i.ID = id
}

func (i *Item) GetCreatedAt() time.Time {
	return i.CreatedAt
}

func (i *Item) SetCreatedAt(t time.Time) {
	i.CreatedAt = t
}

func (i *Item) GetUpdatedAt() time.Time {
	return i.UpdatedAt
}

func (i *Item) SetUpdatedAt(t time.Time) {
	i.UpdatedAt = t
}
