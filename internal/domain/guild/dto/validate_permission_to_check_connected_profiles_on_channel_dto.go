package dto

type ValidatePermissionToCheckConnectedProfilesOnChannelDTO struct {
	ProfileID int `json:"profile_id"`
	UserID    int `json:"user_id"`
	GuildID   int `json:"guild_id"`
}
