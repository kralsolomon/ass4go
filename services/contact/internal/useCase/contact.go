package useCase

import (
	"context"
	"time"

	"advanced.microservices/services/contact/internal/domain"
)

type contactUsecase struct {
	contactRepo    domain.ContactRepository
	contextTimeout time.Duration
}

// Create implements domain.ContactUseCase
func (uc *contactUsecase) Create(contact *domain.Contact) error {
	ctx, cancel := context.WithTimeout(context.Background(), uc.contextTimeout)
	defer cancel()

	return uc.contactRepo.Create(contact, ctx)
}

// Delete implements domain.ContactUseCase
func (uc *contactUsecase) Delete(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), uc.contextTimeout)
	defer cancel()

	return uc.contactRepo.Delete(id, ctx)
}

// GetByID implements domain.ContactUseCase
func (uc *contactUsecase) GetByID(id int64) (*domain.Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), uc.contextTimeout)
	defer cancel()

	return uc.contactRepo.GetByID(id, ctx)
}

// Update implements domain.ContactUseCase
func (uc *contactUsecase) Update(contact *domain.Contact) error {
	ctx, cancel := context.WithTimeout(context.Background(), uc.contextTimeout)
	defer cancel()

	return uc.contactRepo.Update(contact, ctx)
}

func NewContactUsecase(c domain.ContactRepository, timeout time.Duration) domain.ContactUseCase {
	return &contactUsecase{
		contactRepo:    c,
		contextTimeout: timeout,
	}
}
