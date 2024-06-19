package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *UserService) Register(dto dto.RegisterDTO) error {
	if _, err := s.userRepo.FindByEmail(dto.Email); err == nil {
		return validation.NewConflictErr("user", "email")
	}

	user, err := entity.NewUser(
		dto.FirstName,
		dto.LastName,
		dto.Email,
		dto.Password,
	)
	if err != nil {
		return err
	}

	passwordHash, err := s.crypto.GenerateHash(user.PasswordHash)
	if err != nil {
		return validation.NewInternalErr()
	}
	user.PasswordHash = passwordHash

	if err := s.userRepo.Store(user); err != nil {
		return validation.NewInternalErr()
	}

	return nil
}
