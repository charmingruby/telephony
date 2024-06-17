package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *UserService) GetProfileByID(id int) (*entity.UserProfile, error) {
	profile, err := s.profileRepo.FindByUserID(id)
	if err != nil {
		return nil, validation.NewNotFoundErr("user profile")
	}

	return profile, nil
}
