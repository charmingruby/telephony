package integration

import (
	"encoding/json"
	"net/http"

	"github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
)

func (s *Suite) Test_RegisterEndpoint() {
	s.Run("it should be able to register a new user", func() {
		payload := endpoint.RegisterRequest{
			FirstName: "dummy firstname",
			LastName:  "dummy lastname",
			Email:     "dummy@email.com",
			Password:  "password123",
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		res, err := http.Post(s.Route("/v1/auth/register"), contentType, writeBody(body))
		s.NoError(err)
		defer res.Body.Close()

		s.Equal(http.StatusCreated, res.StatusCode)

		var data endpoint.Response
		err = parsePayload[endpoint.Response](&data, res.Body)
		s.NoError(err)

		s.Equal("user created successfully", data.Message)
		s.Equal(http.StatusCreated, data.Code)
	})

	s.Run("it should be not able to register a new user with validation errors", func() {
		payload := endpoint.RegisterRequest{
			FirstName: "dummy firstname",
			LastName:  "dummy lastname",
			Email:     "invalid email",
			Password:  "password123",
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		res, err := http.Post(s.Route("/v1/auth/register"), contentType, writeBody(body))
		s.NoError(err)
		defer res.Body.Close()

		s.Equal(http.StatusUnprocessableEntity, res.StatusCode)

		var data endpoint.Response
		err = parsePayload[endpoint.Response](&data, res.Body)
		s.NoError(err)

		s.Equal("invalid email format", data.Message)
		s.Equal(http.StatusUnprocessableEntity, data.Code)
	})

	s.Run("it should be not able to register a new user with payload errors", func() {
		payload := endpoint.RegisterRequest{}
		body, err := json.Marshal(payload)
		s.NoError(err)

		res, err := http.Post(s.Route("/v1/auth/register"), contentType, writeBody(body))
		s.NoError(err)
		defer res.Body.Close()

		s.Equal(http.StatusBadRequest, res.StatusCode)

		var data endpoint.Response
		err = parsePayload[endpoint.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(http.StatusBadRequest, data.Code)
	})
}
