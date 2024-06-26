package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *UserService) CreateProfile(dto dto.CreateProfileDTO) error {
	_, err := s.userRepo.FindByID(dto.UserID)
	if err != nil {
		return validation.NewNotFoundErr("user")
	}

	_, err = s.profileRepo.FindByDisplayName(dto.DisplayName)
	if err == nil {
		return validation.NewConflictErr("user profile", "display_name")
	}

	profile, err := entity.NewUserProfile(dto.DisplayName, dto.Bio, dto.UserID)
	if err != nil {
		return err
	}

	if _, err := s.profileRepo.Store(profile); err != nil {
		return validation.NewInternalErr()
	}

	return nil
}
