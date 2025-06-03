package request

import (
	"errors"
)

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (r *RegisterUserRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username cannot be empty")
	}
	if len(r.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if r.Password == "" {
		return errors.New("password cannot be empty")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	if r.Email == "" {
		return errors.New("email cannot be empty")
	}
	if !isValidEmail(r.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginUserRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username cannot be empty")
	}
	if r.Password == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

type UpdateUserRoleRequest struct {
	RoleName string `json:"role"`
}

func (r *UpdateUserRoleRequest) Validate() error {
	if r.RoleName == "" {
		return errors.New("role name cannot be empty")
	}
	// if r.RoleName != "user" && r.RoleName != "admin" && r.RoleName != "author" {
	// 	return errors.New("invalid role name: must be 'user', 'admin' or 'author'")
	// }
	return nil
}

func isValidEmail(email string) bool {
	return len(email) > 3 && (contains(email, "@") && contains(email, "."))
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
