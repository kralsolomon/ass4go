package repository

import (
	"context"
	"database/sql"
	"errors"

	"advanced.microservices/services/contact/internal/domain"
)

type SQLContactRepository struct {
	DB *sql.DB
}

// Create implements domain.ContactRepository
func (repository *SQLContactRepository) Create(contact *domain.Contact, ctx context.Context) error {
	query := `
		INSERT INTO contacts (full_name, phone)
		VALUES ($1, $2)
		RETURNING id, created_at`
	args := []any{contact.FullName, contact.Phone}
	err := repository.DB.QueryRowContext(ctx, query, args...).Scan(&contact.ID, &contact.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Delete implements domain.ContactRepository
func (repository *SQLContactRepository) Delete(id int64, ctx context.Context) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM contacts
		WHERE id = $1`

	result, err := repository.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// GetByID implements domain.ContactRepository
func (repository *SQLContactRepository) GetByID(id int64, ctx context.Context) (*domain.Contact, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, full_name, phone, version
		FROM contacts
		WHERE id = $1`

	var contact domain.Contact

	err := repository.DB.QueryRowContext(ctx, query, id).Scan(
		&contact.ID,
		&contact.FullName,
		&contact.Phone,
		&contact.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &contact, nil
}

// Update implements domain.ContactRepository
func (repository *SQLContactRepository) Update(contact *domain.Contact, ctx context.Context) error {
	query := `
		UPDATE contacts
		SET full_name = $1, phone = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version`

	args := []any{
		contact.FullName,
		contact.Phone,
		contact.ID,
		contact.Version,
	}

	err := repository.DB.QueryRowContext(ctx, query, args...).Scan(&contact.Version)
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

func NewContactRepository(conn *sql.DB) domain.ContactRepository {
	return &SQLContactRepository{conn}
}
