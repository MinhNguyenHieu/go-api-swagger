package store

import (
	"context"
	"database/sql"
	"fmt"

	"external-backend-go/db/sqlc"
	"external-backend-go/internal/model"
)

type SessionStore interface {
	CreateSession(ctx context.Context, session *model.Session) (*model.Session, error)
	GetSessionByID(ctx context.Context, id string) (*model.Session, error)
	UpdateSession(ctx context.Context, session *model.Session) (*model.Session, error)
	DeleteSession(ctx context.Context, id string) error
	DeleteExpiredSessions(ctx context.Context, lastActivity int32) error
}

type sessionStore struct {
	*BaseRepository
	queries *sqlc.Queries
}

func NewSessionStore(db *sql.DB, queries *sqlc.Queries, baseRepo *BaseRepository) SessionStore {
	return &sessionStore{BaseRepository: baseRepo, queries: queries}
}

func (s *sessionStore) CreateSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	params := sqlc.CreateSessionParams{
		ID:           session.ID,
		UserID:       session.UserID,
		IpAddress:    session.IpAddress.ToSQLNullString(),
		UserAgent:    session.UserAgent.ToSQLNullString(),
		Payload:      session.Payload,
		LastActivity: session.LastActivity,
	}
	createdSession, err := s.queries.CreateSession(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create session in DB: %w", err)
	}
	return &model.Session{
		ID:           createdSession.ID,
		UserID:       createdSession.UserID,
		IpAddress:    model.FromSQLNullString(createdSession.IpAddress),
		UserAgent:    model.FromSQLNullString(createdSession.UserAgent),
		Payload:      createdSession.Payload,
		LastActivity: createdSession.LastActivity,
	}, nil
}

func (s *sessionStore) GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	dbSession, err := s.queries.GetSessionByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get session by ID from DB: %w", err)
	}
	return &model.Session{
		ID:           dbSession.ID,
		UserID:       dbSession.UserID,
		IpAddress:    model.FromSQLNullString(dbSession.IpAddress),
		UserAgent:    model.FromSQLNullString(dbSession.UserAgent),
		Payload:      dbSession.Payload,
		LastActivity: dbSession.LastActivity,
	}, nil
}

func (s *sessionStore) UpdateSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	params := sqlc.UpdateSessionParams{
		ID:           session.ID,
		UserID:       session.UserID,
		IpAddress:    session.IpAddress.ToSQLNullString(),
		UserAgent:    session.UserAgent.ToSQLNullString(),
		Payload:      session.Payload,
		LastActivity: session.LastActivity,
	}
	updatedSession, err := s.queries.UpdateSession(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to update session in DB: %w", err)
	}
	return &model.Session{
		ID:           updatedSession.ID,
		UserID:       updatedSession.UserID,
		IpAddress:    model.FromSQLNullString(updatedSession.IpAddress),
		UserAgent:    model.FromSQLNullString(updatedSession.UserAgent),
		Payload:      updatedSession.Payload,
		LastActivity: updatedSession.LastActivity,
	}, nil
}

func (s *sessionStore) DeleteSession(ctx context.Context, id string) error {
	err := s.queries.DeleteSession(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to delete session from DB: %w", err)
	}
	return nil
}

func (s *sessionStore) DeleteExpiredSessions(ctx context.Context, lastActivity int32) error {
	err := s.queries.DeleteExpiredSessions(ctx, lastActivity)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions from DB: %w", err)
	}
	return nil
}
