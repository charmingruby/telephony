package dto

import "github.com/charmingruby/telephony/internal/core"

type FetchGuildChannelsDTO struct {
	GuildID    string `json:"guild_id"`
	Pagination core.PaginationParams
}
