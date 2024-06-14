package entity

import (
	"time"

	"github.com/charmingruby/telephony/internal/core"
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

	// validation

	return &profile, nil
}

type UserProfile struct {
	ID               int        `json:"id" db:"id"`
	DisplayName      string     `json:"display_name" db:"display_name"`
	Bio              string     `json:"bio" db:"bio"`
	GuildsQuantity   int        `json:"guilds_quantity" db:"guilds_quantity"`
	MessagesQuantity int        `json:"messages_quantity" db:"messages_quantity"`
	UserID           int        `json:"user_id" db:"user_id"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at" db:"deleted_at"`
}
