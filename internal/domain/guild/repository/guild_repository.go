package repository

import (
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
)

type GuildRepository interface {
	Store(g *entity.Guild) (int, error)
	FindByID(id int) (*entity.Guild, error)
	FindByName(name string) (*entity.Guild, error)
	ListAvailables(page int) ([]entity.Guild, error)
	Delete(g *entity.Guild) error
}
