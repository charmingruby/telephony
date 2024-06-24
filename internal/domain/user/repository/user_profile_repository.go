package repository

import "github.com/charmingruby/telephony/internal/domain/user/entity"

type UserProfileRepository interface {
	Store(u *entity.UserProfile) (int, error)
	FindByID(id int) (*entity.UserProfile, error)
	FindByUserID(userID int) (*entity.UserProfile, error)
	FindByDisplayName(displayName string) (*entity.UserProfile, error)
	UpdateGuildsQuantity(id int, quantity int) error
	UpdateMessagesQuantity(id int, quantity int) error
}
