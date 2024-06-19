package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_CredentialsAuth() {
	fakeCrypto := s.userService.crypto
	dummyEmail := "dummy@email.com"
	dummyPassword := "password123"
	passwordHash, _ := fakeCrypto.GenerateHash(dummyPassword)
	dummyUser, _ := entity.NewUser("dummy", "user", dummyEmail, passwordHash)

	s.Run("it should be able to authenticate", func() {
		dummyUser.ID = 2
		err := s.userRepo.Store(dummyUser)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		dto := dto.CredentialsAuthDTO{
			Email:    dummyEmail,
			Password: dummyPassword,
		}

		data, err := s.userService.CredentialsAuth(dto)

		s.NoError(err)
		s.Equal(dummyUser.ID, data.UserID)
	})

	s.Run("it should be not able to authenticate with nonexistent email", func() {
		dto := dto.CredentialsAuthDTO{
			Email:    dummyEmail,
			Password: dummyPassword,
		}

		data, err := s.userService.CredentialsAuth(dto)

		s.Nil(data)
		s.Error(err)
		s.Equal(validation.NewInvalidCredentialsErr().Error(), err.Error())
	})

	s.Run("it should be not able to authenticate with unmatching password", func() {
		dummyUser.ID = 2
		err := s.userRepo.Store(dummyUser)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		dto := dto.CredentialsAuthDTO{
			Email:    dummyEmail,
			Password: dummyPassword + "2",
		}

		data, err := s.userService.CredentialsAuth(dto)

		s.Nil(data)
		s.Error(err)
		s.Equal(validation.NewInvalidCredentialsErr().Error(), err.Error())
	})
}
