package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/user/adapter"
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/domain/user/repository"
)

type UserServiceContract interface {
	Register(dto dto.RegisterDTO) error
	CreateProfile(dto dto.CreateProfileDTO) error
	CredentialsAuth(dto dto.CredentialsAuthDTO) (*credentialsAuthResponse, error)
	GetProfileByID(id int) (*entity.UserProfile, error)
}

func NewUserService(
	userRepo repository.UserRepository,
	profileRepo repository.UserProfileRepository,
	crypto adapter.CryptographyContract,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		crypto:      crypto,
	}
}

type UserService struct {
	userRepo    repository.UserRepository
	profileRepo repository.UserProfileRepository
	crypto      adapter.CryptographyContract
}
