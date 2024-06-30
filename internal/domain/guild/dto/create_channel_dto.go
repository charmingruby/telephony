package dto

type CreateChannelDTO struct {
	Name      string `json:"name"`
	UserID    int    `json:"user_id"`
	ProfileID int    `json:"profile_id"`
	GuildID   int    `json:"guild_id"`
}
