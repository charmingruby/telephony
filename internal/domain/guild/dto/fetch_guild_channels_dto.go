package dto

type FetchGuildChannelsDTO struct {
	UserID    int `json:"user_id"`
	ProfileID int `json:"profile_id"`
	GuildID   int `json:"guild_id"`
}
