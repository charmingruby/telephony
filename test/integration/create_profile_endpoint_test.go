package integration

import (
	"encoding/json"
	"net/http"

	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/infra/transport/common/middleware"
	"github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_CreateProfileEndpoint() {
	s.Run("it should be able to create a profile", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)

		payload := endpoint.CreateProfileRequest{
			DisplayName: "dummy name",
			Bio:         "dummy biography",
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route("/v1/me/profile"), writeBody(body))
		s.NoError(err)
		req.Header.Add("Content-Type", contentType)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusCreated, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal("user profile created successfully", data.Message)
		s.Equal(http.StatusCreated, data.Code)
	})

	s.Run("it should be not able to create a profile when is not authenticated", func() {
		payload := endpoint.RegisterRequest{}
		body, err := json.Marshal(payload)
		s.NoError(err)

		res, err := http.Post(s.Route("/v1/me/profile"), contentType, writeBody(body))
		s.NoError(err)
		defer res.Body.Close()

		s.Equal(http.StatusUnauthorized, res.StatusCode)

		var data middleware.Response
		err = parsePayload[middleware.Response](&data, res.Body)
		s.NoError(err)

		s.Equal("token not found on header: Authorization", data.Message)
		s.Equal(http.StatusUnauthorized, data.Code)
	})

	s.Run("it should be not able to create a profile with invalid payload", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)

		payload := endpoint.CreateProfileRequest{}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route("/v1/me/profile"), writeBody(body))
		s.NoError(err)
		req.Header.Add("Content-Type", contentType)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusBadRequest, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(http.StatusBadRequest, data.Code)
	})

	s.Run("it should be not able to create a profile with conflicting display_name", func() {
		conflictingDisplayName := "dummy name"

		user, err := createSampleUser("dummy1@email.com", s.userRepo)
		s.NoError(err)
		_, err = createSampleUserProfile(user.ID, conflictingDisplayName, s.profileRepo)
		s.NoError(err)

		conflictingUser, err := createSampleUser("dummy2@email.com", s.userRepo)
		s.NoError(err)

		payload := endpoint.CreateProfileRequest{
			DisplayName: conflictingDisplayName,
			Bio:         "dummy biography",
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route("/v1/me/profile"), writeBody(body))
		s.NoError(err)
		req.Header.Add("Content-Type", contentType)
		err = authenticate(s.token, req, conflictingUser.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusConflict, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(validation.NewConflictErr("user profile", "display_name").Error(), data.Message)
		s.Equal(http.StatusConflict, data.Code)

	})

	s.Run("it should be not able to create a profile with domain error", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)

		payload := endpoint.CreateProfileRequest{
			DisplayName: "21",
			Bio:         "dummy biography",
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route("/v1/me/profile"), writeBody(body))
		s.NoError(err)
		req.Header.Add("Content-Type", contentType)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusUnprocessableEntity, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(validation.NewValidationErr(validation.ErrMinLength("displayname", "4")).Error(), data.Message)
		s.Equal(http.StatusUnprocessableEntity, data.Code)
	})
}
