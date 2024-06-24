package repository

import "github.com/charmingruby/telephony/internal/domain/user/entity"

type UserRepository interface {
	Store(u *entity.User) (int, error)
	FindByID(id int) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}
