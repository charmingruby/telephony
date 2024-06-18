package entity

import (
	"time"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewGuild(name, description string, tags []string, ownerID int) (*Guild, error) {
	guild := Guild{
		ID:          core.NewDefaultDomainID(),
		Name:        name,
		Description: description,
		Tags:        tags,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	if err := validation.ValidateStruct(guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

type Guild struct {
	ID          int        `json:"id" validate:"required" db:"id"`
	Name        string     `json:"name" validate:"required,min=1,max=36" db:"name"`
	Description string     `json:"description" validate:"required,min=1,max=255" db:"description"`
	Tags        []string   `json:"tags" validate:"min=1,max=4" db:"tags"`
	OwnerID     int        `json:"owner_id" validate:"required" db:"owner_id"`
	CreatedAt   time.Time  `json:"created_at" validate:"required" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" validate:"required" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
}
