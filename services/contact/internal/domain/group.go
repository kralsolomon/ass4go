package domain

import (
	"context"
	"time"

	"advanced.microservices/pkg/validator"
)

type Group struct {
	ID        int64     `json:"id"`
	GroupName string    `json:"group_name"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version"`
}
type GroupRepository interface {
	Create(Group *Group, ctx context.Context) error
	GetByID(id int64, ctx context.Context) (*Group, error)
	Update(Group *Group, ctx context.Context) error
}

type GroupUseCase interface {
	Create(Group *Group) error
	GetByID(id int64) (*Group, error)
	Update(Group *Group) error
}

func ValidateGroupName(v *validator.Validator, groupName string) {
}

func ValidateGroup(v *validator.Validator, group *Group) {
	v.Check(len(group.GroupName) <= 250, "group name", "must not be longer than 250 characters")
}
