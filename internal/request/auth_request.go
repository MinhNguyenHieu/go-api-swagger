package request

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
	Email    string `json:"email" validate:"required,email"`
}

func (r *RegisterUserRequest) Validate(v *validator.Validate) error {
	if err := v.Struct(r); err != nil {
		return err
	}
	return nil
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r *LoginUserRequest) Validate(v *validator.Validate) error {
	if err := v.Struct(r); err != nil {
		return err
	}
	return nil
}

type UpdateUserRoleRequest struct {
	RoleName string `json:"role" validate:"required,oneof=admin user author"`
}

func (r *UpdateUserRoleRequest) Validate(v *validator.Validate) error {
	if err := v.Struct(r); err != nil {
		return err
	}
	return nil
}

type PasswordResetRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

func (r *PasswordResetRequest) Validate(v *validator.Validate) error {
	if err := v.Struct(r); err != nil {
		return err
	}
	return nil
}
