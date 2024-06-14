package usecase

import "github.com/charmingruby/telephony/internal/domain/user/dto"

type credentialsAuthResponse struct {
	UserID int
}

func (s *UserService) CredentialsAuth(dto dto.CredentialsAuthDTO) (*credentialsAuthResponse, error) {
	return nil, nil
}
