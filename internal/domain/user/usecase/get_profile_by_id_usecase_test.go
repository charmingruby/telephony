package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_GetProfileByID() {
	user, _ := entity.NewUser("dummy name", "dummy lastname", "dummy@email.com", "password123")
	user.ID = 1
	profile, _ := entity.NewUserProfile("dummy_nick", "dummy bio", user.ID)

	s.Run("it should be able to get a profile", func() {
		err := s.userRepo.Store(user)
		s.NoError(err)
		err = s.profileRepo.Store(profile)
		s.NoError(err)

		p, err := s.userService.GetProfileByID(user.ID)

		s.NoError(err)
		s.Equal(profile.ID, p.ID)
		s.Equal(profile.UserID, p.UserID)
	})

	s.Run("it should be not able to get a profile with an invalid UserID", func() {
		p, err := s.userService.GetProfileByID(-1)

		s.Nil(p)
		s.Error(err)
		s.Equal(validation.NewNotFoundErr("user profile").Error(), err.Error())
	})
}
