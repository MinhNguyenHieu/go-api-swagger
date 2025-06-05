package store

import (
	"context"
	"database/sql"
	"fmt"

	"external-backend-go/db/sqlc"
	"external-backend-go/internal/model"
)

type UserStore interface {
	RepositoryInterface[*model.User]

	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByUsername(ctx context.Context, username string) (sqlc.User, error)
	GetUserByID(ctx context.Context, id int32) (sqlc.User, error)
	UpdateUserRole(ctx context.Context, arg sqlc.UpdateUserRoleParams) (sqlc.User, error)
	SoftDeleteUser(ctx context.Context, id int32) (sqlc.User, error)
	RestoreUser(ctx context.Context, id int32) (sqlc.User, error)
	VerifyUserEmail(ctx context.Context, id int32) (sqlc.User, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.User, error)
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
}

type userStore struct {
	*BaseRepository
	queries *sqlc.Queries
}

func NewUserStore(db *sql.DB, queries *sqlc.Queries, baseRepo *BaseRepository) UserStore {
	return &userStore{
		BaseRepository: baseRepo,
		queries:        queries,
	}
}

func (s *userStore) Create(ctx context.Context, user *model.User) (*model.User, error) {
	params := sqlc.CreateUserParams{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		Email:          user.Email,
		RoleID:         user.RoleID,
	}
	createdUser, err := s.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user via generic store: %w", err)
	}
	return &model.User{
		ID:                createdUser.ID,
		Username:          createdUser.Username,
		Email:             createdUser.Email,
		HashedPassword:    createdUser.HashedPassword,
		EmailVerifiedAt:   model.FromSQLNullTime(createdUser.EmailVerifiedAt),
		RoleID:            createdUser.RoleID,
		RememberTokenUUID: model.FromSQLNullString(createdUser.RememberTokenUuid),
		CreatedAt:         createdUser.CreatedAt,
		UpdatedAt:         createdUser.UpdatedAt,
		DeletedAt:         model.FromSQLNullTime(createdUser.DeletedAt),
	}, nil
}

func (s *userStore) GetByID(ctx context.Context, id int32) (*model.User, error) {
	dbUser, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID via generic store: %w", err)
	}
	return &model.User{
		ID:                dbUser.ID,
		Username:          dbUser.Username,
		Email:             dbUser.Email,
		HashedPassword:    dbUser.HashedPassword,
		EmailVerifiedAt:   model.FromSQLNullTime(dbUser.EmailVerifiedAt),
		RoleID:            dbUser.RoleID,
		RememberTokenUUID: model.FromSQLNullString(dbUser.RememberTokenUuid),
		CreatedAt:         dbUser.CreatedAt,
		UpdatedAt:         dbUser.UpdatedAt,
		DeletedAt:         model.FromSQLNullTime(dbUser.DeletedAt),
	}, nil
}

func (s *userStore) Update(ctx context.Context, user *model.User) (*model.User, error) {
	params := sqlc.UpdateUserParams{
		ID:                user.ID,
		Username:          user.Username,
		HashedPassword:    user.HashedPassword,
		Email:             user.Email,
		EmailVerifiedAt:   user.EmailVerifiedAt.ToSQLNullTime(),
		RoleID:            user.RoleID,
		RememberTokenUuid: user.RememberTokenUUID.ToSQLNullString(),
		DeletedAt:         user.DeletedAt.ToSQLNullTime(),
	}
	updatedUser, err := s.queries.UpdateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update user via generic store: %w", err)
	}
	return &model.User{
		ID:                updatedUser.ID,
		Username:          updatedUser.Username,
		Email:             updatedUser.Email,
		HashedPassword:    updatedUser.HashedPassword,
		EmailVerifiedAt:   model.FromSQLNullTime(updatedUser.EmailVerifiedAt),
		RoleID:            updatedUser.RoleID,
		RememberTokenUUID: model.FromSQLNullString(updatedUser.RememberTokenUuid),
		CreatedAt:         updatedUser.CreatedAt,
		UpdatedAt:         updatedUser.UpdatedAt,
		DeletedAt:         model.FromSQLNullTime(updatedUser.DeletedAt),
	}, nil
}

func (s *userStore) Delete(ctx context.Context, id int32) error {
	return s.queries.DeleteUser(ctx, id)
}

func (s *userStore) List(ctx context.Context, offset, limit int32) ([]*model.User, error) {
	params := sqlc.ListUsersParams{
		Offset: offset,
		Limit:  limit,
	}
	dbUsers, err := s.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list users from DB: %w", err)
	}

	var users []*model.User
	for _, dbUser := range dbUsers {
		users = append(users, &model.User{
			ID:                dbUser.ID,
			Username:          dbUser.Username,
			Email:             dbUser.Email,
			HashedPassword:    dbUser.HashedPassword,
			EmailVerifiedAt:   model.FromSQLNullTime(dbUser.EmailVerifiedAt),
			RoleID:            dbUser.RoleID,
			RememberTokenUUID: model.FromSQLNullString(dbUser.RememberTokenUuid),
			CreatedAt:         dbUser.CreatedAt,
			UpdatedAt:         dbUser.UpdatedAt,
			DeletedAt:         model.FromSQLNullTime(dbUser.DeletedAt),
		})
	}
	return users, nil
}

func (s *userStore) Count(ctx context.Context) (int64, error) {
	count, err := s.queries.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count users in DB: %w", err)
	}
	return count, nil
}

func (s *userStore) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := s.queries.CreateUser(ctx, arg)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to create user via store: %w", err)
	}
	return user, nil
}

func (s *userStore) GetUserByUsername(ctx context.Context, username string) (sqlc.User, error) {
	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to get user by username via store: %w", err)
	}
	return user, nil
}

func (s *userStore) GetUserByID(ctx context.Context, id int32) (sqlc.User, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to get user by ID via store: %w", err)
	}
	return user, nil
}

func (s *userStore) UpdateUserRole(ctx context.Context, arg sqlc.UpdateUserRoleParams) (sqlc.User, error) {
	user, err := s.queries.UpdateUserRole(ctx, arg)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to update user role via store: %w", err)
	}
	return user, nil
}

func (s *userStore) SoftDeleteUser(ctx context.Context, id int32) (sqlc.User, error) {
	user, err := s.queries.SoftDeleteUser(ctx, id)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to soft delete user via store: %w", err)
	}
	return user, nil
}

func (s *userStore) RestoreUser(ctx context.Context, id int32) (sqlc.User, error) {
	user, err := s.queries.RestoreUser(ctx, id)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to restore user via store: %w", err)
	}
	return user, nil
}

func (s *userStore) VerifyUserEmail(ctx context.Context, id int32) (sqlc.User, error) {
	user, err := s.queries.VerifyUserEmail(ctx, id)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to verify user email via store: %w", err)
	}
	return user, nil
}

func (s *userStore) GetUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to get user by email via store: %w", err)
	}
	return user, nil
}

func (s *userStore) UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	user, err := s.queries.UpdateUser(ctx, arg)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to update user via store: %w", err)
	}
	return user, nil
}
