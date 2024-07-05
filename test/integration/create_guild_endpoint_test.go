package integration

import (
	"encoding/json"
	"net/http"

	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_CreateGuildEndpoint() {
	dummyName := "dummy name"
	dummyDescription := "dummy description"

	s.Run("it should be able to create a guild", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		payload := endpoint.CreateGuildRequest{
			Name:        dummyName,
			Description: dummyDescription,
			ProfileID:   profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route("/v1/guilds"), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusCreated, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal("guild created successfully", data.Message)
		s.Equal(http.StatusCreated, data.Code)

		modifiedProfile, err := s.profileRepo.FindByID(profile.ID)
		s.NoError(err)
		s.Equal(1, modifiedProfile.GuildsQuantity)
	})

	s.Run("it should be not able to create a guild if profile dont exists", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)

		payload := endpoint.CreateGuildRequest{
			Name:        dummyName,
			Description: dummyDescription,
			ProfileID:   -2,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route("/v1/guilds"), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusNotFound, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(validation.NewNotFoundErr("user_profile").Error(), data.Message)
		s.Equal(http.StatusNotFound, data.Code)

	})

	s.Run("it should be not able to create a guild if user dont exists", func() {
		payload := endpoint.CreateGuildRequest{
			Name:        dummyName,
			Description: dummyDescription,
			ProfileID:   -2,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route("/v1/guilds"), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, -2)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusNotFound, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(validation.NewNotFoundErr("user").Error(), data.Message)
		s.Equal(http.StatusNotFound, data.Code)
	})

	s.Run("it should be not able to create a guild if its not the profile owner", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		ownerUser, err := createSampleUser("dummy2@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(ownerUser.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		payload := endpoint.CreateGuildRequest{
			Name:        dummyName,
			Description: dummyDescription,
			ProfileID:   profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route("/v1/guilds"), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusUnauthorized, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(validation.NewUnauthorizedErr().Error(), data.Message)
		s.Equal(http.StatusUnauthorized, data.Code)
	})

	s.Run("it should be not able to create a guild if name is conflicting", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "guild name", s.guildRepo)
		s.NoError(err)

		payload := endpoint.CreateGuildRequest{
			Name:        guild.Name,
			Description: dummyDescription,
			ProfileID:   profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route("/v1/guilds"), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusConflict, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(validation.NewConflictErr("guild", "name").Error(), data.Message)
		s.Equal(http.StatusConflict, data.Code)
	})

	s.Run("it should be not able to create a guild if have validation errors", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		var name string
		for i := 0; i < 37; i++ {
			name += "i"
		}

		payload := endpoint.CreateGuildRequest{
			Name:        name,
			Description: dummyDescription,
			ProfileID:   profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route("/v1/guilds"), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusUnprocessableEntity, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(http.StatusUnprocessableEntity, data.Code)
	})

	s.Run("it should be not able to create a guild if have invalid params", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		payload := endpoint.CreateGuildRequest{
			Name:        "",
			Description: dummyDescription,
			ProfileID:   profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route("/v1/guilds"), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusBadRequest, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(http.StatusBadRequest, data.Code)
	})
}
