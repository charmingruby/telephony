package dto

type JoinChannelDTO struct {
	ProfileID string `json:"profile_id"`
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}
