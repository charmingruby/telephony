package usecase

import "github.com/charmingruby/telephony/internal/domain/user/repository"

type UserServiceContract interface{}

func NewUserService(
	userRepo repository.UserRepository,
	profileRepo repository.UserProfileRepository,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
	}
}

type UserService struct {
	userRepo    repository.UserRepository
	profileRepo repository.UserProfileRepository
}
