package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_CreateProfile() {
	dummyDisplayName := "dummy name"
	dummyBio := "dummy bio"
	dummyUser, _ := entity.NewUser("dummy", "user", "dummy@email.com", "password123")

	s.Run("it should be able to create a profile with a valid payload", func() {
		dummyUser.ID = 2
		err := s.userRepo.Store(dummyUser)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		dto := dto.CreateProfileDTO{
			DisplayName: dummyDisplayName,
			Bio:         dummyBio,
			UserID:      dummyUser.ID,
		}

		err = s.userService.CreateProfile(dto)

		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))
		s.Equal(dummyDisplayName, s.profileRepo.Items[0].DisplayName)
	})

	s.Run("it should be not able to create a profile with an invalid UserID", func() {
		dto := dto.CreateProfileDTO{
			DisplayName: dummyDisplayName,
			Bio:         dummyBio,
			UserID:      -1,
		}

		err := s.userService.CreateProfile(dto)

		s.Error(err)
		s.Equal(validation.NewNotFoundErr("user").Error(), err.Error())
	})

	s.Run("it should be not able to create a profile with an invalid params", func() {
		dummyUser.ID = 2
		err := s.userRepo.Store(dummyUser)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		dto := dto.CreateProfileDTO{
			DisplayName: "",
			Bio:         dummyBio,
			UserID:      dummyUser.ID,
		}

		err = s.userService.CreateProfile(dto)

		s.Error(err)
		s.Equal(validation.NewValidationErr(validation.ErrRequired("displayname")).Error(), err.Error())
	})
}
