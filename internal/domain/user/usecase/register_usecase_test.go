package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_Register() {
	user, _ := entity.NewUser("dummy name", "dummy lastname", "dummy@email.com", "password123")
	dto := dto.RegisterDTO{
		FirstName: user.FirstName,
		LastName:  user.FirstName,
		Email:     user.Email,
		Password:  user.PasswordHash,
	}

	s.Run("it should be able to register", func() {
		err := s.userService.Register(dto)

		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))
		s.Equal(s.userRepo.Items[0].Email, dto.Email)
	})

	s.Run("it should be not able to register with conflicting email", func() {
		err := s.userRepo.Store(user)
		s.NoError(err)

		err = s.userService.Register(dto)

		s.Error(err)
		s.Equal(validation.NewConflictErr("user", "email").Error(), err.Error())
	})

	s.Run("it should be not able to register with invalid payload", func() {
		dto.FirstName = ""

		err := s.userService.Register(dto)

		s.Error(err)
		s.Equal(validation.NewValidationErr(validation.ErrRequired("firstname")).Error(), err.Error())
	})
}
