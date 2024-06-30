package inmemory

import (
	"github.com/charmingruby/telephony/internal/core"
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

func (r *InMemoryChannelRepository) ListChannelsByGuildID(guildID, page int) ([]entity.Channel, error) {
	channels := []entity.Channel{}

	pageToFilter := page - 1
	startValue := pageToFilter * core.ItemsPerPage()
	endValue := startValue + core.ItemsPerPage()

	if startValue >= len(r.Items) {
		return channels, nil
	}

	if endValue > len(r.Items) {
		endValue = len(r.Items)
	}

	for i := startValue; i < endValue; i++ {
		currentItem := r.Items[i]
		if currentItem.GuildID == guildID {
			channels = append(channels, currentItem)
		}
	}

	return channels, nil
}
