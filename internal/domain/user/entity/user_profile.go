package entity

import (
	"time"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewUserProfile(displayName, bio string, userID int) (*UserProfile, error) {
	profile := UserProfile{
		ID:               core.NewDefaultDomainID(),
		DisplayName:      displayName,
		Bio:              bio,
		GuildsQuantity:   0,
		MessagesQuantity: 0,
		UserID:           userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		DeletedAt:        nil,
	}

	if err := validation.ValidateStruct(profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

type UserProfile struct {
	ID               int        `json:"id" validate:"required" db:"id"`
	DisplayName      string     `json:"display_name" validate:"required,min=4,max=16" db:"display_name"`
	Bio              string     `json:"bio" validate:"required,max=32" db:"bio"`
	GuildsQuantity   int        `json:"guilds_quantity" db:"guilds_quantity"`
	MessagesQuantity int        `json:"messages_quantity" db:"messages_quantity"`
	UserID           int        `json:"user_id" validate:"required" db:"user_id"`
	CreatedAt        time.Time  `json:"created_at" validate:"required" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" validate:"required" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at" db:"deleted_at"`
}
