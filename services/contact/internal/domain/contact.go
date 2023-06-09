package domain

import (
	"context"
	"regexp"
	"strings"
	"time"

	"advanced.microservices/pkg/validator"
)

type Contact struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version"`
}

type ContactRepository interface {
	Create(contact *Contact, ctx context.Context) error
	GetByID(id int64, ctx context.Context) (*Contact, error)
	Update(contact *Contact, ctx context.Context) error
	Delete(id int64, ctx context.Context) error
}

type ContactUseCase interface {
	Create(contact *Contact) error
	GetByID(id int64) (*Contact, error)
	Update(contact *Contact) error
	Delete(id int64) error
}

func ValidateContact(v *validator.Validator, contact *Contact) {
	v.Check(len(strings.Split(contact.FullName, " ")) == 3, "full name", "full name must contain 3 parts")
	v.Check(validator.Matches(contact.Phone, regexp.MustCompile(`[0-9\[\]\(\\)\+\-]`)), "email", "must be a valid phone number")
}
