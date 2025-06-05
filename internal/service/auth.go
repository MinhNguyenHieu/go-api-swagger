package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"external-backend-go/db/sqlc"
	"external-backend-go/internal/auth"
	"external-backend-go/internal/email"
	"external-backend-go/internal/model"
	"external-backend-go/internal/store"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidRoleName      = errors.New("invalid role name specified")
	ErrIncorrectPassword    = errors.New("incorrect username or password")
	ErrInvalidToken         = errors.New("invalid or expired token")
	ErrEmailAlreadyVerified = errors.New("email already verified")
	ErrUserAlreadyExists    = errors.New("user with this username or email already exists")
)

type AuthService struct {
	UserStore               store.UserStore
	RoleStore               store.RoleStore
	SessionStore            store.SessionStore
	PasswordResetTokenStore store.PasswordResetTokenStore
	JWTSecret               string
	EmailSender             email.EmailSender
}

func NewAuthService(userStore store.UserStore, roleStore store.RoleStore, sessionStore store.SessionStore, passwordResetTokenStore store.PasswordResetTokenStore, jwtSecret string, emailSender email.EmailSender) *AuthService {
	return &AuthService{
		UserStore:               userStore,
		RoleStore:               roleStore,
		SessionStore:            sessionStore,
		PasswordResetTokenStore: passwordResetTokenStore,
		JWTSecret:               jwtSecret,
		EmailSender:             emailSender,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, username, password, email, roleName string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	role, err := s.RoleStore.GetRoleByName(ctx, roleName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidRoleName // Role not found
		}
		return fmt.Errorf("failed to get role by name: %w", err)
	}

	arg := sqlc.CreateUserParams{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Email:          email,
		RoleID:         role.ID,
	}

	_, err = s.UserStore.CreateUser(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *AuthService) LoginUser(ctx context.Context, username, password, ipAddress, userAgent string) (string, string, error) {
	dbUser, err := s.UserStore.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrUserNotFound
		}
		return "", "", fmt.Errorf("failed to get user by username: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.HashedPassword), []byte(password))
	if err != nil {
		return "", "", ErrIncorrectPassword
	}

	role, err := s.RoleStore.GetByID(ctx, dbUser.RoleID)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user role: %w", err)
	}

	// Create session
	sessionPayload, err := json.Marshal(map[string]string{"username": dbUser.Username, "role": role.Name})
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal session payload: %w", err)
	}

	session := &model.Session{
		ID:           uuid.New().String(),
		UserID:       sql.NullInt32{Int32: dbUser.ID, Valid: true},
		IpAddress:    model.NullString{String: ipAddress, Valid: ipAddress != ""},
		UserAgent:    model.NullString{String: userAgent, Valid: userAgent != ""},
		Payload:      string(sessionPayload),
		LastActivity: int32(time.Now().Unix()),
	}

	_, err = s.SessionStore.CreateSession(ctx, session)
	if err != nil {
		return "", "", fmt.Errorf("failed to create session: %w", err)
	}

	tokenString, err := auth.GenerateToken(dbUser.ID, dbUser.Username, role.Name, s.JWTSecret)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, role.Name, nil
}

// VerifyEmail verifies the user's email address.
func (s *AuthService) VerifyEmail(ctx context.Context, userID int32, token string) error {
	user, err := s.UserStore.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user for email verification: %w", err)
	}

	if user.EmailVerifiedAt.Valid {
		return ErrEmailAlreadyVerified
	}

	if token == "" {
		return ErrInvalidToken
	}

	_, err = s.UserStore.VerifyUserEmail(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to verify user email in store: %w", err)
	}

	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	_, err := s.UserStore.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user by email for password reset: %w", err)
	}

	token := uuid.New().String()
	expiryTime := time.Now().Add(15 * time.Minute)

	existingToken, err := s.PasswordResetTokenStore.GetPasswordResetToken(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to check existing password reset token: %w", err)
	}

	if existingToken != nil {
		err = s.PasswordResetTokenStore.DeletePasswordResetToken(ctx, email)
		if err != nil {
			return fmt.Errorf("failed to delete old password reset token: %w", err)
		}
	}
	_, err = s.PasswordResetTokenStore.CreatePasswordResetToken(ctx, email, token, expiryTime)
	if err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}

	// Construct email body with reset link
	resetLink := fmt.Sprintf("http://localhost:8080/reset-password?email=%s&token=%s", email, token)
	body := fmt.Sprintf(`
		<p>Hello,</p>
		<p>You have requested to reset your password. Please click the link below to reset it:</p>
		<p><a href="%s">%s</a></p>
		<p>This link will expire in 15 minutes.</p>
		<p>If you did not request a password reset, please ignore this email.</p>
	`, resetLink, resetLink)

	// Send email
	err = s.EmailSender.SendEmail(email, "Password Reset Request", body)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, email, token, newPassword string) error {
	resetToken, err := s.PasswordResetTokenStore.GetPasswordResetToken(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidToken
		}
		return fmt.Errorf("failed to get password reset token: %w", err)
	}

	if resetToken.Token != token || !resetToken.CreatedAt.Valid || time.Now().After(resetToken.CreatedAt.Time.Add(15*time.Minute)) {
		return ErrInvalidToken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	user, err := s.UserStore.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user by email for password reset: %w", err)
	}

	arg := sqlc.UpdateUserParams{
		ID:                user.ID,
		Username:          user.Username,
		HashedPassword:    string(hashedPassword),
		Email:             user.Email,
		EmailVerifiedAt:   user.EmailVerifiedAt,
		RoleID:            user.RoleID,
		RememberTokenUuid: user.RememberTokenUuid,
		DeletedAt:         user.DeletedAt,
	}

	_, err = s.UserStore.UpdateUser(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	err = s.PasswordResetTokenStore.DeletePasswordResetToken(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to delete password reset token: %w", err)
	}

	return nil
}

func (s *AuthService) UpdateUserRole(ctx context.Context, userID int32, newRoleName string) (sqlc.User, error) {
	_, err := s.UserStore.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, ErrUserNotFound
		}
		return sqlc.User{}, fmt.Errorf("failed to retrieve user: %w", err)
	}

	newRole, err := s.RoleStore.GetRoleByName(ctx, newRoleName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, ErrInvalidRoleName
		}
		return sqlc.User{}, fmt.Errorf("failed to get new role by name: %w", err)
	}

	updatedUser, err := s.UserStore.UpdateUserRole(ctx, sqlc.UpdateUserRoleParams{
		ID:     userID,
		RoleID: newRole.ID,
	})
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to update user role: %w", err)
	}
	return updatedUser, nil
}
