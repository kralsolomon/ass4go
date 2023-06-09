package repository

import (
	"context"
	"database/sql"
	"errors"

	"advanced.microservices/services/contact/internal/domain"
)

type SQLGroupRepository struct {
	DB *sql.DB
}

// Create implements domain.GroupRepository
func (repository *SQLGroupRepository) Create(group *domain.Group, ctx context.Context) error {
	query := `
		INSERT INTO groups (group_name)
		VALUES ($1)
		RETURNING id, created_at`
	args := []any{group.GroupName}

	err := repository.DB.QueryRowContext(ctx, query, args...).Scan(&group.ID, &group.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetByID implements domain.GroupRepository
func (repository *SQLGroupRepository) GetByID(id int64, ctx context.Context) (*domain.Group, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, full_name, phone, version
		FROM groups
		WHERE id = $1`

	var group domain.Group

	err := repository.DB.QueryRowContext(ctx, query, id).Scan(
		&group.ID,
		&group.GroupName,
		&group.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &group, nil
}

// Update implements domain.GroupRepository
func (repository *SQLGroupRepository) Update(group *domain.Group, ctx context.Context) error {
	query := `
		UPDATE groups
		SET group_name = $1, version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version`

	args := []any{
		group.GroupName,
		group.ID,
		group.Version,
	}

	err := repository.DB.QueryRowContext(ctx, query, args...).Scan(&group.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func NewGroupRepository(conn *sql.DB) domain.GroupRepository {
	return &SQLGroupRepository{conn}
}
