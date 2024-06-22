package repository

import (
	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
)

type GuildRepository interface {
	Store(g *entity.Guild) error
	ListAvailables(pagination core.PaginationParams)
	Delete(g *entity.Guild) error
}
