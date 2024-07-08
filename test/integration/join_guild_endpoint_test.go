package integration

import (
	"encoding/json"
	"fmt"
	"net/http"

	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_JoinGuildEndpoint() {
	s.Run("it should be able to join a guild sucessfully", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)

		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		guild, err := createSampleGuild(profile.ID, "dummy guild name", s.guildRepo)
		s.NoError(err)

		payload := endpoint.JoinGuildRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/join", guild.ID)
		req, err := http.NewRequest(http.MethodPost, s.Route(rawRoute), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusOK, res.StatusCode)

		var data connhelper.Response
		err = parsePayload[connhelper.Response](&data, res.Body)
		s.NoError(err)

		s.Equal("user joined to the guild successfully", data.Message)
		s.Equal(http.StatusOK, data.Code)
	})

	s.Run("it should be not able to join a guild if is already a member", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)

		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		guild, err := createSampleGuild(profile.ID, "dummy guild name", s.guildRepo)
		s.NoError(err)

		_, err = createSampleMember(user.ID, profile.ID, guild.ID, s.guildMemberRepo)
		s.NoError(err)

		payload := endpoint.JoinGuildRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/join", guild.ID)
		req, err := http.NewRequest(http.MethodPost, s.Route(rawRoute), writeBody(body))
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

		s.Equal("user profile is already a member", data.Message)
		s.Equal(http.StatusBadRequest, data.Code)
	})

	s.Run("it should be not able to join a guild if profile dont exists", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy guild name", s.guildRepo)
		s.NoError(err)

		payload := endpoint.JoinGuildRequest{
			ProfileID: -2,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/join", guild.ID)
		req, err := http.NewRequest(http.MethodPost, s.Route(rawRoute), writeBody(body))
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

	s.Run("it should be not able to join a guild if user dont exists", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy guild name", s.guildRepo)
		s.NoError(err)

		payload := endpoint.JoinGuildRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/join", guild.ID)
		req, err := http.NewRequest(http.MethodPost, s.Route(rawRoute), writeBody(body))
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
		guild, err := createSampleGuild(profile.ID, "dummy guild name", s.guildRepo)
		s.NoError(err)

		payload := endpoint.JoinGuildRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/join", guild.ID)
		req, err := http.NewRequest(http.MethodPost, s.Route(rawRoute), writeBody(body))
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
}
