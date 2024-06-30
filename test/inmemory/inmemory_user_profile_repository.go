package inmemory

import (
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewInMemoryUserProfileRepository() *InMemoryUserProfileRepository {
	return &InMemoryUserProfileRepository{
		Items: []entity.UserProfile{},
	}
}

type InMemoryUserProfileRepository struct {
	Items []entity.UserProfile
}

func (r *InMemoryUserProfileRepository) Store(e *entity.UserProfile) (int, error) {
	r.Items = append(r.Items, *e)
	return e.ID, nil
}

func (r *InMemoryUserProfileRepository) FindByID(id int) (*entity.UserProfile, error) {
	for _, e := range r.Items {
		if e.ID == id {
			return &e, nil
		}
	}

	return nil, validation.NewNotFoundErr("user profile")
}

func (r *InMemoryUserProfileRepository) FindByUserID(userID int) (*entity.UserProfile, error) {
	for _, e := range r.Items {
		if e.UserID == userID {
			return &e, nil
		}
	}

	return nil, validation.NewNotFoundErr("user profile")
}

func (r *InMemoryUserProfileRepository) FindByDisplayName(displayName string) (*entity.UserProfile, error) {
	for _, e := range r.Items {
		if e.DisplayName == displayName {
			return &e, nil
		}
	}

	return nil, validation.NewNotFoundErr("user profile")
}

func (r *InMemoryUserProfileRepository) UpdateGuildsQuantity(id int, quantity int) error {
	for idx, e := range r.Items {
		if e.ID == id {
			r.Items[idx].GuildsQuantity = e.GuildsQuantity + quantity
			return nil
		}
	}

	return validation.NewNotFoundErr("user profile")
}

func (r *InMemoryUserProfileRepository) UpdateMessagesQuantity(id int, quantity int) error {
	for idx, e := range r.Items {
		if e.ID == id {
			r.Items[idx].MessagesQuantity = e.MessagesQuantity + quantity
			return nil
		}
	}

	return validation.NewNotFoundErr("user profile")
}
