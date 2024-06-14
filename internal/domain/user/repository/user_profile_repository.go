package repository

import "github.com/charmingruby/telephony/internal/domain/user/entity"

type UserProfileRepository interface {
	Store(u *entity.UserProfile) error
	FindByUserID(userID int) (*entity.UserProfile, error)
	FindByDisplayName(displayName string) (*entity.UserProfile, error)
}
