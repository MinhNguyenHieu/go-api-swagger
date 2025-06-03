package store

import (
	"context"
	"database/sql"
	"fmt"

	"external-backend-go/db/sqlc"
	"external-backend-go/internal/model"
)

type RoleStore interface {
	RepositoryInterface[*model.Role]

	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
}

type roleStore struct {
	*BaseRepository
	queries *sqlc.Queries
}

func NewRoleStore(db *sql.DB, queries *sqlc.Queries, baseRepo *BaseRepository) RoleStore {
	return &roleStore{BaseRepository: baseRepo, queries: queries}
}

func (s *roleStore) Create(ctx context.Context, role *model.Role) (*model.Role, error) {
	params := sqlc.CreateRoleParams{
		Name:        role.Name,
		Description: role.Description.ToSQLNullString(),
	}
	createdRole, err := s.queries.CreateRole(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create role in DB: %w", err)
	}
	return &model.Role{
		ID:          createdRole.ID,
		Name:        createdRole.Name,
		Description: model.FromSQLNullString(createdRole.Description),
		CreatedAt:   createdRole.CreatedAt,
		UpdatedAt:   createdRole.UpdatedAt,
	}, nil
}

func (s *roleStore) GetByID(ctx context.Context, id int32) (*model.Role, error) {
	dbRole, err := s.queries.GetRoleByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get role by ID from DB: %w", err)
	}
	return &model.Role{
		ID:          dbRole.ID,
		Name:        dbRole.Name,
		Description: model.FromSQLNullString(dbRole.Description),
		CreatedAt:   dbRole.CreatedAt,
		UpdatedAt:   dbRole.UpdatedAt,
	}, nil
}

func (s *roleStore) Update(ctx context.Context, role *model.Role) (*model.Role, error) {
	params := sqlc.UpdateRoleParams{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description.ToSQLNullString(),
	}
	updatedRole, err := s.queries.UpdateRole(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to update role in DB: %w", err)
	}
	return &model.Role{
		ID:          updatedRole.ID,
		Name:        updatedRole.Name,
		Description: model.FromSQLNullString(updatedRole.Description),
		CreatedAt:   updatedRole.CreatedAt,
		UpdatedAt:   updatedRole.UpdatedAt,
	}, nil
}

func (s *roleStore) Delete(ctx context.Context, id int32) error {
	err := s.queries.DeleteRole(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to delete role from DB: %w", err)
	}
	return nil
}

func (s *roleStore) List(ctx context.Context, offset, limit int32) ([]*model.Role, error) {
	params := sqlc.ListRolesParams{
		Offset: offset,
		Limit:  limit,
	}
	dbRoles, err := s.queries.ListRoles(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles from DB: %w", err)
	}

	var roles []*model.Role
	for _, dbRole := range dbRoles {
		roles = append(roles, &model.Role{
			ID:          dbRole.ID,
			Name:        dbRole.Name,
			Description: model.FromSQLNullString(dbRole.Description),
			CreatedAt:   dbRole.CreatedAt,
			UpdatedAt:   dbRole.UpdatedAt,
		})
	}
	return roles, nil
}

func (s *roleStore) Count(ctx context.Context) (int64, error) {
	count, err := s.queries.CountRoles(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count roles in DB: %w", err)
	}
	return count, nil
}

func (s *roleStore) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	dbRole, err := s.queries.GetRoleByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get role by name from DB: %w", err)
	}
	return &model.Role{
		ID:          dbRole.ID,
		Name:        dbRole.Name,
		Description: model.FromSQLNullString(dbRole.Description),
		CreatedAt:   dbRole.CreatedAt,
		UpdatedAt:   dbRole.UpdatedAt,
	}, nil
}
