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

// Define sentinel errors for specific conditions
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidRoleName   = errors.New("invalid role name specified")
	ErrIncorrectPassword = errors.New("incorrect username or password")
)

type AuthService struct {
	UserStore    store.UserStore
	RoleStore    store.RoleStore
	SessionStore store.SessionStore
	JWTSecret    string
	EmailSender  email.EmailSender
}

func NewAuthService(userStore store.UserStore, roleStore store.RoleStore, sessionStore store.SessionStore, jwtSecret string, emailSender email.EmailSender) *AuthService {
	return &AuthService{UserStore: userStore, RoleStore: roleStore, SessionStore: sessionStore, JWTSecret: jwtSecret, EmailSender: emailSender}
}

func (s *AuthService) RegisterUser(ctx context.Context, username, password, email, roleName string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	role, err := s.RoleStore.GetRoleByName(ctx, roleName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidRoleName
		}
		return fmt.Errorf("failed to get role: %w", err)
	}

	params := sqlc.CreateUserParams{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Email:          email,
		RoleID:         role.ID,
	}

	_, err = s.UserStore.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	go func() {
		if s.EmailSender != nil {
			welcomeSubject := "Welcome to our application!"
			welcomeBody := fmt.Sprintf("Hello %s,<br><br>Thank you for registering an account with us. Your role is: %s. We are excited to welcome you!<br><br>Regards,<br>Your Application Team", username, role.Name)
			if err := s.EmailSender.SendEmail(email, welcomeSubject, welcomeBody); err != nil {
				fmt.Printf("Error sending welcome email to %s: %v\n", email, err)
			}
		} else {
			fmt.Println("EmailSender not initialized, cannot send welcome email.")
		}
	}()

	return nil
}

func (s *AuthService) LoginUser(ctx context.Context, username, password, ipAddress, userAgent string) (string, string, error) {
	dbUser, err := s.UserStore.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrIncorrectPassword
		}
		return "", "", fmt.Errorf("failed to retrieve user from DB: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.HashedPassword), []byte(password))
	if err != nil {
		return "", "", ErrIncorrectPassword
	}

	role, err := s.RoleStore.GetByID(ctx, dbUser.RoleID)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user role: %w", err)
	}

	sessionID := uuid.New().String()

	sessionPayload, err := json.Marshal(map[string]interface{}{
		"user_id":    dbUser.ID,
		"username":   dbUser.Username,
		"role":       role.Name,
		"login_time": time.Now().Unix(),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal session payload: %w", err)
	}

	session := &model.Session{
		ID:           sessionID,
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

func (s *AuthService) UpdateUserRole(ctx context.Context, userID int32, newRoleName string) (sqlc.User, error) {
	existingUser, err := s.UserStore.GetUserByID(ctx, userID)
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
		return sqlc.User{}, fmt.Errorf("failed to get new role: %w", err)
	}

	if existingUser.RoleID == newRole.ID {
		return existingUser, nil
	}

	params := sqlc.UpdateUserRoleParams{
		ID:     userID,
		RoleID: newRole.ID,
	}

	updatedUser, err := s.UserStore.UpdateUserRole(ctx, params)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to update user role: %w", err)
	}

	return updatedUser, nil
}
