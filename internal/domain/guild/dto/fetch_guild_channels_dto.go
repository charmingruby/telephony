package dto

import "github.com/charmingruby/telephony/internal/core"

type FetchGuildChannelsDTO struct {
	GuildID    int `json:"guild_id"`
	Pagination core.PaginationParams
}
