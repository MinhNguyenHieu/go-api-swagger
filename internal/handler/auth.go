package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"external-backend-go/internal/logger"
	"external-backend-go/internal/request"
	"external-backend-go/internal/service"
	"external-backend-go/internal/utility"
)

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type AuthHandler struct {
	AuthService *service.AuthService
	Logger      *logger.Logger
	Validator   *validator.Validate
}

func NewAuthHandler(authService *service.AuthService, logger *logger.Logger, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{AuthService: authService, Logger: logger, Validator: validator}
}

// @Summary Register new user
// @Description Creates a new user account with username, password, and email. Defaults to 'user' role.
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body request.RegisterUserRequest true "User registration details"
// @Success 201 {object} map[string]string "message: Registration successful!"
// @Failure 400 {object} map[string]string "message: Invalid request data"
// @Failure 500 {object} map[string]string "message: Could not register user. Username or email might already exist."
// @Router /register [post]
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req request.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid request data"), h.Logger)
		return
	}

	if err := req.Validate(h.Validator); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			utility.BadRequestResponse(w, r, fmt.Errorf("Validation failed: %s", ve.Error()), h.Logger)
			return
		}
		utility.BadRequestResponse(w, r, err, h.Logger)
		return
	}

	err := h.AuthService.RegisterUser(r.Context(), req.Username, req.Password, req.Email, "user")
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			utility.InternalServerError(w, r, fmt.Errorf("Username or email already exists"), h.Logger)
		} else {
			utility.InternalServerError(w, r, err, h.Logger)
		}
		return
	}

	utility.JSONResponse(w, http.StatusCreated, map[string]string{"message": "Registration successful!"})
}

// @Summary Login user
// @Description Logs in a user and returns a JWT token.
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body request.LoginUserRequest true "User login credentials"
// @Success 200 {object} LoginResponse "Successful login"
// @Failure 400 {object} map[string]string "message: Invalid request data"
// @Failure 401 {object} map[string]string "message: Invalid username or password"
// @Failure 500 {object} map[string]string "message: Internal server error"
// @Router /login [post]
func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req request.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid request data"), h.Logger)
		return
	}

	if err := req.Validate(h.Validator); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			utility.BadRequestResponse(w, r, fmt.Errorf("Validation failed: %s", ve.Error()), h.Logger)
			return
		}
		utility.BadRequestResponse(w, r, err, h.Logger)
		return
	}

	token, role, err := h.AuthService.LoginUser(r.Context(), req.Username, req.Password, r.RemoteAddr, r.UserAgent())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrIncorrectPassword) {
			utility.UnauthorizedErrorResponse(w, r, fmt.Errorf("Invalid username or password"), h.Logger)
		} else {
			utility.InternalServerError(w, r, err, h.Logger)
		}
		return
	}

	utility.JSONResponse(w, http.StatusOK, LoginResponse{Token: token, Role: role})
}

// @Summary Verify user email
// @Description Verifies a user's email address using a provided token.
// @Tags authentication
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]string "message: Email verified successfully!"
// @Failure 400 {object} map[string]string "message: Invalid verification link or token"
// @Failure 500 {object} map[string]string "message: Internal server error"
// @Router /verify-email [get]
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid User ID format"), h.Logger)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		utility.BadRequestResponse(w, r, fmt.Errorf("Verification token is missing"), h.Logger)
		return
	}

	err = h.AuthService.VerifyEmail(r.Context(), int32(userID), token)
	if err != nil {
		utility.InternalServerError(w, r, err, h.Logger)
		return
	}

	utility.JSONResponse(w, http.StatusOK, map[string]string{"message": "Email verified successfully!"})
}

// @Summary Protected Endpoint
// @Description This is a sample protected endpoint accessible only with a valid JWT.
// @Tags example
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} map[string]string "message: Access granted!"
// @Failure 401 {object} map[string]string "message: Authentication token required / Invalid token"
// @Router /protected [get]
func (h *AuthHandler) ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {

	utility.JSONResponse(w, http.StatusOK, map[string]string{"message": "Access granted to protected endpoint!"})
}

// @Summary Request password reset
// @Description Sends a password reset email to the user.
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body map[string]string true "email: User's email address"
// @Success 200 {object} map[string]string "message: Password reset email sent."
// @Failure 400 {object} map[string]string "message: Invalid email format / Email not found"
// @Failure 500 {object} map[string]string "message: Failed to send password reset email."
// @Router /forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid request data"), h.Logger)
		return
	}

	if err := h.Validator.Struct(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			utility.BadRequestResponse(w, r, fmt.Errorf("Validation failed: %s", ve.Error()), h.Logger)
			return
		}
		utility.BadRequestResponse(w, r, err, h.Logger)
		return
	}

	if err := h.AuthService.ForgotPassword(r.Context(), req.Email); err != nil {
		utility.InternalServerError(w, r, err, h.Logger)
		return
	}

	utility.JSONResponse(w, http.StatusOK, map[string]string{"message": "Password reset email sent."})
}

// @Summary Reset password
// @Description Resets user's password using a valid token.
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body map[string]string true "email: User's email, token: Reset token, new_password: New password"
// @Success 200 {object} map[string]string "message: Password reset successfully!"
// @Failure 400 {object} map[string]string "message: Invalid request data / Invalid or expired token / Passwords do not match criteria"
// @Failure 500 {object} map[string]string "message: Internal server error"
// @Router /reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req request.PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid request data"), h.Logger)
		return
	}

	if err := req.Validate(h.Validator); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			utility.BadRequestResponse(w, r, fmt.Errorf("Validation failed: %s", ve.Error()), h.Logger)
			return
		}
		utility.BadRequestResponse(w, r, err, h.Logger)
		return
	}

	err := h.AuthService.ResetPassword(r.Context(), req.Email, req.Token, req.NewPassword)
	if err != nil {
		utility.InternalServerError(w, r, err, h.Logger)
		return
	}

	utility.JSONResponse(w, http.StatusOK, map[string]string{"message": "Password reset successfully!"})
}

// @Summary Update user role (Admin only)
// @Description Updates the role of a specific user. Requires JWT authentication and 'admin' role.
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Param request body request.UpdateUserRoleRequest true "New role for the user"
// @Success 200 {object} model.User "Updated user details"
// @Failure 400 {object} map[string]string "message: Invalid request data / Invalid User ID format"
// @Failure 401 {object} map[string]string "message: Authentication token required / Invalid token"
// @Failure 403 {object} map[string]string "message: You do not have permission to access this resource."
// @Failure 404 {object} map[string]string "message: User not found"
// @Failure 500 {object} map[string]string "message: Internal server error"
// @Router /admin/users/{id}/role [put]
func (h *AuthHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid User ID format"), h.Logger)
		return
	}

	var req request.UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid request data"), h.Logger)
		return
	}

	if err := req.Validate(h.Validator); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			utility.BadRequestResponse(w, r, fmt.Errorf("Validation failed: %s", ve.Error()), h.Logger)
			return
		}
		utility.BadRequestResponse(w, r, err, h.Logger)
		return
	}

	updatedUser, err := h.AuthService.UpdateUserRole(r.Context(), int32(userID), req.RoleName)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrInvalidRoleName) {
			utility.NotFoundResponse(w, r, h.Logger)
		} else {
			utility.InternalServerError(w, r, err, h.Logger)
		}
		return
	}

	role, err := h.AuthService.RoleStore.GetByID(r.Context(), updatedUser.RoleID)
	if err != nil {
		h.Logger.Error("Failed to fetch role for updated user %d: %v", updatedUser.ID, err)
		utility.InternalServerError(w, r, fmt.Errorf("Failed to retrieve role information"), h.Logger)
		return
	}

	response := map[string]interface{}{
		"id":                updatedUser.ID,
		"username":          updatedUser.Username,
		"email":             updatedUser.Email,
		"emailVerifiedAt":   updatedUser.EmailVerifiedAt,
		"roleId":            updatedUser.RoleID,
		"roleName":          role.Name,
		"rememberTokenUuid": updatedUser.RememberTokenUuid,
		"createdAt":         updatedUser.CreatedAt,
		"updatedAt":         updatedUser.UpdatedAt,
		"deletedAt":         updatedUser.DeletedAt,
	}
	utility.JSONResponse(w, http.StatusOK, response)
}

// @Summary Protected with Basic Auth Endpoint
// @Description This is a sample protected endpoint accessible only with Basic Authentication.
// @Tags example
// @Security BasicAuth
// @Produce json
// @Success 200 {object} map[string]string "message: Basic Auth access granted!"
// @Failure 401 {object} map[string]string "message: Basic authentication failed"
// @Router /basic-auth/protected [get]
func (h *AuthHandler) ProtectedWithBasicAuth(w http.ResponseWriter, r *http.Request) {
	utility.JSONResponse(w, http.StatusOK, map[string]string{"message": "Basic Auth access granted!"})
}
