package repository

import "github.com/charmingruby/telephony/internal/domain/user/entity"

type UserRepository interface {
	Store(u *entity.User) error
	FindByUUID(uuid string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}
