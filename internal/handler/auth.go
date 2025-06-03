package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
}

func NewAuthHandler(authService *service.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{AuthService: authService, Logger: logger}
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

	if err := req.Validate(); err != nil {
		utility.BadRequestResponse(w, r, err, h.Logger)
		return
	}

	err := h.AuthService.RegisterUser(r.Context(), req.Username, req.Password, req.Email, "user")
	if err != nil {
		utility.InternalServerError(w, r, err, h.Logger)
		return
	}

	utility.JSONResponse(w, http.StatusCreated, map[string]string{"message": "Registration successful!"})
}

// @Summary Log in user and get JWT token
// @Description Authenticates user with username and password, then returns a JWT token and user role.
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body request.LoginUserRequest true "User login credentials"
// @Success 200 {object} LoginResponse "token: JWT token, role: User's role"
// @Failure 400 {object} map[string]string "message: Invalid request data"
// @Failure 401 {object} map[string]string "message: Incorrect username or password"
// @Failure 500 {object} map[string]string "message: Internal server error"
// @Router /login [post]
func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req request.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utility.BadRequestResponse(w, r, fmt.Errorf("Invalid request data"), h.Logger)
		return
	}

	if err := req.Validate(); err != nil {
		utility.BadRequestResponse(w, r, err, h.Logger)
		return
	}

	ipAddress := getClientIP(r)
	userAgent := r.UserAgent()

	tokenString, userRole, err := h.AuthService.LoginUser(r.Context(), req.Username, req.Password, ipAddress, userAgent)
	if err != nil {
		if errors.Is(err, service.ErrIncorrectPassword) {
			utility.UnauthorizedErrorResponse(w, r, err, h.Logger)
		} else {
			utility.InternalServerError(w, r, err, h.Logger)
		}
		return
	}

	utility.JSONResponse(w, http.StatusOK, LoginResponse{Token: tokenString, Role: userRole})
}

// falling back to r.RemoteAddr.
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}
	ip, _, found := strings.Cut(r.RemoteAddr, ":")
	if !found {
		return r.RemoteAddr
	}
	return ip
}

// @Summary Access protected endpoint
// @Description This endpoint requires JWT authentication.
// @Tags protected
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "message: You have successfully accessed the protected endpoint!"
// @Failure 401 {object} map[string]string "message: Authentication token required / Invalid token"
// @Router /protected [get]
func (h *AuthHandler) ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	utility.JSONResponse(w, http.StatusOK, map[string]string{"message": "You have successfully accessed the protected endpoint!"})
}

// @Summary Update user role by ID
// @Description Allows an administrator to change a user's role. Requires 'admin' role.
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Param request body request.UpdateUserRoleRequest true "New role details"
// @Success 200 {object} map[string]string "message: User role updated successfully!"
// @Failure 400 {object} map[string]string "message: Invalid request data / Invalid User ID"
// @Failure 401 {object} map[string]string "message: Authentication token required / Invalid token"
// @Failure 403 {object} map[string]string "message: Forbidden (requires admin role)"
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

	if err := req.Validate(); err != nil {
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

	// Fetch the role name from the RoleStore using the updatedUser.RoleID
	role, err := h.AuthService.RoleStore.GetByID(r.Context(), updatedUser.RoleID)
	if err != nil {
		utility.InternalServerError(w, r, fmt.Errorf("Failed to retrieve updated role name: %w", err), h.Logger)
		return
	}

	utility.JSONResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("User %s role updated to %s successfully!", updatedUser.Username, role.Name)})
}
