package dto

type JoinChannelDTO struct {
	UserID    int `json:"user_id"`
	ProfileID int `json:"profile_id"`
	GuildID   int `json:"guild_id"`
	ChannelID int `json:"channel_id"`
}

type JoinChannelResponseDTO struct {
	DisplayName string `json:"display_name"`
}
