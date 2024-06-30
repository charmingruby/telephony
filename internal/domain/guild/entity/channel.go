package entity

import (
	"time"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewChannel(name string, guildID, profileID int) (*Channel, error) {
	channel := Channel{
		ID:               core.NewDefaultDomainID(),
		Name:             name,
		MessagesQuantity: 0,
		GuildID:          guildID,
		OwnerID:          profileID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		DeletedAt:        nil,
	}

	if err := validation.ValidateStruct(channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

type Channel struct {
	ID               int        `json:"id" validate:"required" db:"id"`
	Name             string     `json:"name" validate:"required,min=1,max=36" db:"name"`
	MessagesQuantity int        `json:"messages_quantity" db:"messages_quantity"`
	GuildID          int        `json:"guild_id" validate:"required" db:"guild_id"`
	OwnerID          int        `json:"owner_id" validate:"required" db:"owner_id"`
	CreatedAt        time.Time  `json:"created_at" validate:"required" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" validate:"required" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at" db:"deleted_at"`
}
