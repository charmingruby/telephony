package dto

type DeleteMessageDTO struct {
	MessageID string `json:"message_id"`
	SenderID  int    `json:"sender_id"`
	ChannelID int    `json:"channel_id"`
}
