package inmemory

import (
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		Items: []entity.User{},
	}
}

type InMemoryUserRepository struct {
	Items []entity.User
}

func (r *InMemoryUserRepository) Store(e *entity.User) error {
	r.Items = append(r.Items, *e)
	return nil
}

func (r *InMemoryUserRepository) FindByUUID(uuid string) (*entity.User, error) {
	for _, e := range r.Items {
		if e.UUID == uuid {
			return &e, nil
		}
	}

	return nil, validation.NewNotFoundErr("user")
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*entity.User, error) {
	for _, e := range r.Items {
		if e.Email == email {
			return &e, nil
		}
	}

	return nil, validation.NewNotFoundErr("user")
}
