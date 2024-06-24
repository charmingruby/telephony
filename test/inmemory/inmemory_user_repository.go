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

func (r *InMemoryUserRepository) Store(e *entity.User) (int, error) {
	r.Items = append(r.Items, *e)
	return e.ID, nil
}

func (r *InMemoryUserRepository) FindByID(id int) (*entity.User, error) {
	for _, e := range r.Items {
		if e.ID == id {
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
