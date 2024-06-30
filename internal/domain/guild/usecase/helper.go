package usecase

import (
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *GuildService) userProfileValidation(profileID, userID int) error {
	userExists := s.userClient.UserExists(userID)
	if !userExists {
		return validation.NewNotFoundErr("user")
	}

	profileExists := s.userClient.ProfileExists(profileID)
	if !profileExists {
		return validation.NewNotFoundErr("user_profile")
	}

	isTheProfileOwner := s.userClient.IsTheProfileOwner(userID, profileID)
	if !isTheProfileOwner {
		return validation.NewUnauthorizedErr()
	}

	return nil
}
