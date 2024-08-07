package repository

import "github.com/charmingruby/telephony/internal/domain/guild/entity"

type ChannelRepository interface {
	Store(c *entity.Channel) (int, error)
	FindByName(guildID int, name string) (*entity.Channel, error)
	FindByID(channelID, guildID int) (*entity.Channel, error)
	ListChannelsByGuildID(guildID int) ([]entity.Channel, error)
	ListAllChannels() ([]entity.Channel, error)
}
