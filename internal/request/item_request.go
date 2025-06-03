package request

import "errors"

type CreateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *CreateItemRequest) Validate() error {
	if r.Name == "" {
		return errors.New("item name cannot be empty")
	}
	if len(r.Name) < 3 {
		return errors.New("item name must be at least 3 characters long")
	}
	return nil
}

type UpdateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *UpdateItemRequest) Validate() error {
	if r.Name == "" {
		return errors.New("item name cannot be empty")
	}
	if len(r.Name) < 3 {
		return errors.New("item name must be at least 3 characters long")
	}
	return nil
}
