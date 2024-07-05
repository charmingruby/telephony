package integration

import (
	"encoding/json"
	"fmt"
	"net/http"

	connhelper "github.com/charmingruby/telephony/internal/infra/transport/common/conn_helper"
	"github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_CreateChannelEndpoint() {
	dummyName := "dummy name"

	s.Run("it should be able to create a channel", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy name", s.guildRepo)
		s.NoError(err)
		_, err = createSampleMember(profile.ID, user.ID, guild.ID, s.guildMemberRepo)
		s.NoError(err)

		payload := endpoint.CreateChannelRequest{
			Name:      dummyName,
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route(fmt.Sprintf("/v1/guilds/%d/channels", guild.ID)), writeBody(body))
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

		s.Equal("channel created successfully", data.Message)
		s.Equal(http.StatusCreated, data.Code)
	})

	s.Run("it should be not able to create a channel if its not a member", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy name", s.guildRepo)
		s.NoError(err)

		payload := endpoint.CreateChannelRequest{
			Name:      dummyName,
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route(fmt.Sprintf("/v1/guilds/%d/channels", guild.ID)), writeBody(body))
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

	s.Run("it should be not able to create a channel if profile dont exists", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy name", s.guildRepo)
		s.NoError(err)

		payload := endpoint.CreateChannelRequest{
			Name:      dummyName,
			ProfileID: -2,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route(fmt.Sprintf("/v1/guilds/%d/channels", guild.ID)), writeBody(body))
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

	s.Run("it should be not able to create a channel if user dont exists", func() {
		payload := endpoint.CreateChannelRequest{
			Name:      dummyName,
			ProfileID: -2,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route(fmt.Sprintf("/v1/guilds/%d/channels", 1)), writeBody(body))
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

	s.Run("it should be not able to create a channel if guild dont exists", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		payload := endpoint.CreateChannelRequest{
			Name:      dummyName,
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route(fmt.Sprintf("/v1/guilds/%d/channels", -1)), writeBody(body))
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

		s.Equal(validation.NewNotFoundErr("guild").Error(), data.Message)
		s.Equal(http.StatusNotFound, data.Code)
	})

	s.Run("it should be not able to create a guild if its not the profile owner", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		ownerUser, err := createSampleUser("dummy2@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(ownerUser.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy name", s.guildRepo)
		s.NoError(err)

		payload := endpoint.CreateChannelRequest{
			Name:      dummyName,
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route(fmt.Sprintf("/v1/guilds/%d/channels", guild.ID)), writeBody(body))
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

	s.Run("it should be not able to create a channel if guild already have a channel with this name", func() {
		conflictingName := "dummy name"

		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "guild name", s.guildRepo)
		s.NoError(err)
		_, err = createSampleMember(profile.ID, user.ID, guild.ID, s.guildMemberRepo)
		s.NoError(err)
		_, err = createSampleChannel(conflictingName, profile.ID, guild.ID, s.channelRepo)
		s.NoError(err)

		payload := endpoint.CreateChannelRequest{
			Name:      conflictingName,
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route(fmt.Sprintf("/v1/guilds/%d/channels", guild.ID)), writeBody(body))
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

		s.Equal(validation.NewConflictErr("channel", "name").Error(), data.Message)
		s.Equal(http.StatusConflict, data.Code)
	})

	s.Run("it should be not able to create a guild if have validation errors", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy name", s.guildRepo)
		s.NoError(err)
		_, err = createSampleMember(profile.ID, user.ID, guild.ID, s.guildMemberRepo)
		s.NoError(err)

		payload := endpoint.CreateChannelRequest{
			Name:      "",
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, s.Route(fmt.Sprintf("/v1/guilds/%d/channels", guild.ID)), writeBody(body))
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
