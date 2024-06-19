package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/user/dto"
	"github.com/charmingruby/telephony/internal/validation"
)

type credentialsAuthResponse struct {
	UserID int
}

func (s *UserService) CredentialsAuth(dto dto.CredentialsAuthDTO) (*credentialsAuthResponse, error) {
	user, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil {
		return nil, validation.NewInvalidCredentialsErr()
	}

	isCredentialsValid := s.crypto.ValidateHash(dto.Password, user.PasswordHash)
	if !isCredentialsValid {
		return nil, validation.NewInvalidCredentialsErr()
	}

	// fields to be saved on token
	return &credentialsAuthResponse{
		UserID: user.ID, // sub
	}, nil
}
