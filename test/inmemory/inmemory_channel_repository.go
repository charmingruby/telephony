package inmemory

import (
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewInMemoryChannelRepository() *InMemoryChannelRepository {
	return &InMemoryChannelRepository{
		Items: []entity.Channel{},
	}
}

type InMemoryChannelRepository struct {
	Items []entity.Channel
}

func (r *InMemoryChannelRepository) Store(e *entity.Channel) (int, error) {
	r.Items = append(r.Items, *e)
	return e.ID, nil
}

func (r *InMemoryChannelRepository) FindByName(guildID int, name string) (*entity.Channel, error) {
	for _, e := range r.Items {
		if e.Name == name && e.GuildID == guildID {
			return &e, nil
		}
	}

	return nil, validation.NewNotFoundErr("channel")
}

func (r *InMemoryChannelRepository) FindByID(channelID, guildID int) (*entity.Channel, error) {
	for _, e := range r.Items {
		if e.ID == channelID && e.GuildID == guildID {
			return &e, nil
		}
	}

	return nil, validation.NewNotFoundErr("channel")
}

func (r *InMemoryChannelRepository) ListChannelsByGuildID(guildID int) ([]entity.Channel, error) {
	channels := []entity.Channel{}

	for i := 0; i < len(r.Items); i++ {
		currentItem := r.Items[i]
		if currentItem.GuildID == guildID {
			channels = append(channels, currentItem)
		}
	}

	return channels, nil
}

func (r *InMemoryChannelRepository) ListAllChannels() ([]entity.Channel, error) {
	return r.Items, nil
}
