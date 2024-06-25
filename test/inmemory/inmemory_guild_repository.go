package inmemory

import (
	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewInMemoryGuildRepository() *InMemoryGuildRepository {
	return &InMemoryGuildRepository{
		Items: []entity.Guild{},
	}
}

type InMemoryGuildRepository struct {
	Items []entity.Guild
}

func (r *InMemoryGuildRepository) Store(e *entity.Guild) (int, error) {
	r.Items = append(r.Items, *e)
	return e.ID, nil
}

func (r *InMemoryGuildRepository) FindByID(id int) (*entity.Guild, error) {
	for _, e := range r.Items {
		if e.ID == id {
			return &e, nil
		}
	}

	return nil, validation.NewNotFoundErr("guild")
}

func (r *InMemoryGuildRepository) FindByName(name string) (*entity.Guild, error) {
	for _, e := range r.Items {
		if e.Name == name {
			return &e, nil
		}
	}

	return nil, validation.NewNotFoundErr("guild")
}

func (r *InMemoryGuildRepository) ListAvailables(page int) ([]entity.Guild, error) {
	guilds := []entity.Guild{}

	pageToFilter := page - 1
	startValue := pageToFilter * core.ItemsPerPage()
	endValue := startValue + core.ItemsPerPage()

	for i := startValue; i < endValue; i++ {
		currentItem := r.Items[i]
		if currentItem.DeletedAt == nil {
			guilds = append(guilds, currentItem)
		}
	}

	return guilds, nil
}

func (r *InMemoryGuildRepository) Delete(e *entity.Guild) error {
	var idx int

	for _, g := range r.Items {
		if g.ID == e.ID {
			idx = e.ID
		}
	}

	r.Items = append(r.Items[:idx], r.Items[idx+1])

	return nil
}
