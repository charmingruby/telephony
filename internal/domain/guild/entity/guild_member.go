package entity

import (
	"time"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewGuildMember(profileID, userID, guildID int) (*GuildMember, error) {
	guildMember := GuildMember{
		ID:        core.NewDefaultDomainID(),
		ProfileID: profileID,
		UserID:    userID,
		GuildID:   guildID,
		IsActive:  true,
		JoinedAt:  time.Now(),
	}

	if err := validation.ValidateStruct(guildMember); err != nil {
		return nil, err
	}

	return &guildMember, nil
}

type GuildMember struct {
	ID        int       `json:"id" validate:"required" db:"id"`
	ProfileID int       `json:"profile_id" validate:"required" db:"profile_id"`
	UserID    int       `json:"user_id" validate:"required" db:"user_id"`
	GuildID   int       `json:"guild_id" validate:"required" db:"guild_id"`
	IsActive  bool      `json:"is_active" validate:"required" db:"is_active"`
	JoinedAt  time.Time `json:"joined_at" validate:"required" db:"joined_at"`
}
