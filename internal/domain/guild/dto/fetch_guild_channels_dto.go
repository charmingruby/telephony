package dto

import "github.com/charmingruby/telephony/internal/core"

type FetchGuildChannelsDTO struct {
	UserID     int `json:"user_id"`
	ProfileID  int `json:"profile_id"`
	GuildID    int `json:"guild_id"`
	Pagination core.PaginationParams
}
