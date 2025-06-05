package request

import "github.com/go-playground/validator/v10"

type CreateItemRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

func (r *CreateItemRequest) Validate(v *validator.Validate) error {
	if err := v.Struct(r); err != nil {
		return err
	}
	return nil
}

type UpdateItemRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

func (r *UpdateItemRequest) Validate(v *validator.Validate) error {
	if err := v.Struct(r); err != nil {
		return err
	}
	return nil
}
