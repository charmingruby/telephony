package integration

import (
	"encoding/json"
	"net/http"

	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/infra/security/cryptography"
	"github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
)

func (s *Suite) Test_CredentialsAuthEndpoint() {
	password := "password123"
	crypto := cryptography.NewCryptography()
	passwordHash, _ := crypto.GenerateHash(password)
	user, _ := entity.NewUser(
		"dummy name",
		"dummy lastname",
		"dummy@email.com",
		password,
	)
	user.ID = 1
	user.PasswordHash = passwordHash

	s.Run("it should be able to authenticate with valid credentials", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)

		payload := endpoint.CredentialsAuthRequest{
			Email:    user.Email,
			Password: password,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		res, err := http.Post(s.Route("/v1/auth/login"), contentType, writeBody(body))
		s.NoError(err)

		s.Equal(http.StatusOK, res.StatusCode)

		var data endpoint.CredentialsAuthResponse
		err = parsePayload[endpoint.CredentialsAuthResponse](&data, res.Body)
		s.NoError(err)

		s.Equal("user authenticated successfully", data.Message)
		s.Equal(http.StatusOK, data.Code)

		// validate token payload
		token := data.Data.AccessToken
		tokenPayload, err := s.token.ValidateToken(token)
		s.NoError(err)
		s.Equal(user.ID, tokenPayload.UserID)
	})

	s.Run("it should be not able to authenticate with invalid email", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)

		payload := endpoint.CredentialsAuthRequest{
			Email:    "invalid@email.com",
			Password: password,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		res, err := http.Post(s.Route("/v1/auth/login"), contentType, writeBody(body))
		s.NoError(err)

		s.Equal(http.StatusUnauthorized, res.StatusCode)

		var data endpoint.Response
		err = parsePayload[endpoint.Response](&data, res.Body)
		s.NoError(err)

		s.Equal("invalid credentials", data.Message)
		s.Equal(http.StatusUnauthorized, data.Code)
	})

	s.Run("it should be not able to authenticate with invalid password", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)

		payload := endpoint.CredentialsAuthRequest{
			Email:    user.Email,
			Password: "invalid password",
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		res, err := http.Post(s.Route("/v1/auth/login"), contentType, writeBody(body))
		s.NoError(err)

		s.Equal(http.StatusUnauthorized, res.StatusCode)

		var data endpoint.Response
		err = parsePayload[endpoint.Response](&data, res.Body)
		s.NoError(err)

		s.Equal("invalid credentials", data.Message)
		s.Equal(http.StatusUnauthorized, data.Code)
	})

	s.Run("it should be not able to authenticate with invalid payload", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)

		payload := endpoint.CredentialsAuthRequest{}
		body, err := json.Marshal(payload)
		s.NoError(err)

		res, err := http.Post(s.Route("/v1/auth/login"), contentType, writeBody(body))
		s.NoError(err)

		s.Equal(http.StatusBadRequest, res.StatusCode)

		var data endpoint.Response
		err = parsePayload[endpoint.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(http.StatusBadRequest, data.Code)
	})

}
