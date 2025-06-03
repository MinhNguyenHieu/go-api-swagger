package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"external-backend-go/db/sqlc"
	"external-backend-go/internal/model"
)

type PasswordResetTokenStore interface {
	CreatePasswordResetToken(ctx context.Context, email, token string, createdAt time.Time) (*model.PasswordResetToken, error)
	GetPasswordResetToken(ctx context.Context, email string) (*model.PasswordResetToken, error)
	DeletePasswordResetToken(ctx context.Context, email string) error
}

type passwordResetTokenStore struct {
	*BaseRepository
	queries *sqlc.Queries
}

func NewPasswordResetTokenStore(db *sql.DB, queries *sqlc.Queries, baseRepo *BaseRepository) PasswordResetTokenStore {
	return &passwordResetTokenStore{BaseRepository: baseRepo, queries: queries}
}

func (s *passwordResetTokenStore) CreatePasswordResetToken(ctx context.Context, email, token string, createdAt time.Time) (*model.PasswordResetToken, error) {
	params := sqlc.CreatePasswordResetTokenParams{
		Email:     email,
		Token:     token,
		CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
	}
	createdToken, err := s.queries.CreatePasswordResetToken(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create password reset token in DB: %w", err)
	}
	return &model.PasswordResetToken{
		Email:     createdToken.Email,
		Token:     createdToken.Token,
		CreatedAt: model.FromSQLNullTime(createdToken.CreatedAt),
	}, nil
}

func (s *passwordResetTokenStore) GetPasswordResetToken(ctx context.Context, email string) (*model.PasswordResetToken, error) {
	dbToken, err := s.queries.GetPasswordResetToken(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get password reset token from DB: %w", err)
	}
	return &model.PasswordResetToken{
		Email:     dbToken.Email,
		Token:     dbToken.Token,
		CreatedAt: model.FromSQLNullTime(dbToken.CreatedAt),
	}, nil
}

func (s *passwordResetTokenStore) DeletePasswordResetToken(ctx context.Context, email string) error {
	err := s.queries.DeletePasswordResetToken(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to delete password reset token from DB: %w", err)
	}
	return nil
}
