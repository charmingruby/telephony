package entity

import (
	"time"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewChannelMessage(content string, channelID, senderID int) (*ChannelMessage, error) {
	channelMessage := ChannelMessage{
		ID:        core.NewDefaultDomainID(),
		Content:   content,
		ChannelID: channelID,
		SenderID:  senderID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	if err := validation.ValidateStruct(channelMessage); err != nil {
		return nil, err
	}

	return &channelMessage, nil
}

type ChannelMessage struct {
	ID        int        `json:"id" validate:"required" db:"id"`
	Content   string     `json:"content" validate:"required,min=1,max=255" db:"content"`
	SenderID  int        `json:"sender_id" validate:"required" db:"sender_id"`
	ChannelID int        `json:"channel_id" validate:"required" db:"channel_id"`
	CreatedAt time.Time  `json:"created_at" validate:"required" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" validate:"required" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}
