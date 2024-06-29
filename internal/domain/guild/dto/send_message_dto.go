package dto

type SendMessageDTO struct {
	Content   string `json:"content"`
	SenderID  int    `json:"sender_id"`
	ChannelID int    `json:"channel_id"`
}
