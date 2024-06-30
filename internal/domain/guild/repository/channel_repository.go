package repository

import "github.com/charmingruby/telephony/internal/domain/guild/entity"

type ChannelRepository interface {
	Store(c *entity.Channel) (int, error)
	FindByName(guildID int, name string) (*entity.Channel, error)
	ListChannelsByGuildID(guildID, page int) ([]entity.Channel, error)
}
