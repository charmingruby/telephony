package integration

import (
	"net/http"

	"github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
	"github.com/charmingruby/telephony/internal/infra/transport/rest/middleware"
)

func (s *Suite) Test_MeEndpoint() {
	s.Run("it should be able to get an user profile", func() {
		user, profile, err := createSampleUser(s.userRepo, s.profileRepo)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, s.Route("/v1/me"), nil)
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusOK, res.StatusCode)

		var data endpoint.MeResponse
		err = parsePayload[endpoint.MeResponse](&data, res.Body)
		s.NoError(err)

		s.Equal("profile found", data.Message)
		s.Equal(profile.DisplayName, data.Data.DisplayName)
		s.Equal(user.ID, data.Data.UserID)
		s.Equal(http.StatusOK, data.Code)
	})

	s.Run("it should be not able to get an user profile when is unauthenticated", func() {
		res, err := http.Get(s.Route("/v1/me"))
		s.NoError(err)
		defer res.Body.Close()

		s.Equal(http.StatusUnauthorized, res.StatusCode)

		var data middleware.Response
		err = parsePayload[middleware.Response](&data, res.Body)
		s.NoError(err)

		s.Equal("token not found on header: Authorization", data.Message)
		s.Equal(http.StatusUnauthorized, data.Code)
	})
}
