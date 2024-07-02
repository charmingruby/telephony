package integration

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_FetchGuildChannelsEndpoint() {
	dummyChannelName := "dummy name"

	s.Run("it should be able fetch guild channels by page", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guildToBeFiltedBy, err := createSampleGuild(profile.ID, "dummy guild name 1", s.guildRepo)
		s.NoError(err)
		anotherGuild, err := createSampleGuild(profile.ID, "dummy guild name 2", s.guildRepo)
		s.NoError(err)
		_, err = createSampleMember(profile.ID, user.ID, guildToBeFiltedBy.ID, s.guildMemberRepo)
		s.NoError(err)

		totalChannels := 4
		createdChannels := []entity.Channel{}
		for i := 0; i < totalChannels; i++ {
			if i%2 == 0 {
				channel, err := createSampleChannel(
					fmt.Sprintf("%s-%d", dummyChannelName, i),
					profile.ID,
					guildToBeFiltedBy.ID,
					s.channelRepo,
				)
				s.NoError(err)
				createdChannels = append(createdChannels, *channel)
			} else {
				_, err := createSampleChannel(
					fmt.Sprintf("%s-%d", dummyChannelName, i),
					profile.ID,
					anotherGuild.ID,
					s.channelRepo,
				)
				s.NoError(err)
			}
		}

		payload := endpoint.FetchGuildChannelsRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/channels?page=1", guildToBeFiltedBy.ID)
		req, err := http.NewRequest(http.MethodGet, s.Route(rawRoute), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusOK, res.StatusCode)

		var data endpoint.FetchGuildChannelsResponse
		err = parsePayload[endpoint.FetchGuildChannelsResponse](&data, res.Body)
		s.NoError(err)

		s.Equal("guild channels fetched", data.Message)
		s.Equal(totalChannels/2, len(data.Data))
		s.Equal(http.StatusOK, data.Code)

		for idx, g := range data.Data {
			s.Equal(createdChannels[idx].ID, g.ID)
			s.Equal(createdChannels[idx].Name, g.Name)
			s.Equal(guildToBeFiltedBy.ID, g.GuildID)
		}
	})

	s.Run("it should be not able fetch guild channels if its not a guild member", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy guild name 1", s.guildRepo)
		s.NoError(err)

		payload := endpoint.FetchGuildChannelsRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/channels?page=1", guild.ID)
		req, err := http.NewRequest(http.MethodGet, s.Route(rawRoute), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusUnauthorized, res.StatusCode)

		var data endpoint.FetchGuildChannelsResponse
		err = parsePayload[endpoint.FetchGuildChannelsResponse](&data, res.Body)
		s.NoError(err)

		s.Equal(http.StatusUnauthorized, data.Code)
		s.Equal(validation.NewUnauthorizedErr().Error(), data.Message)
	})

	s.Run("it should be able fetch guild channels even if is passing an empty page", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy guild name", s.guildRepo)
		s.NoError(err)
		_, err = createSampleMember(profile.ID, user.ID, guild.ID, s.guildMemberRepo)
		s.NoError(err)

		payload := endpoint.FetchGuildChannelsRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/channels?page=1", guild.ID)
		req, err := http.NewRequest(http.MethodGet, s.Route(rawRoute), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusOK, res.StatusCode)

		var data endpoint.FetchGuildChannelsResponse
		err = parsePayload[endpoint.FetchGuildChannelsResponse](&data, res.Body)
		s.NoError(err)

		s.Equal("guild channels fetched", data.Message)
		s.Equal(0, len(data.Data))
		s.Equal(http.StatusOK, data.Code)
	})

	s.Run("it should be able fetch guild channels in a different page from the default", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy guild name", s.guildRepo)
		s.NoError(err)
		_, err = createSampleMember(profile.ID, user.ID, guild.ID, s.guildMemberRepo)
		s.NoError(err)

		totalChannels := core.ItemsPerPage() + 4
		createdChannels := []entity.Channel{}
		for i := 0; i < totalChannels; i++ {
			channel, err := createSampleChannel(
				fmt.Sprintf("%s-%d", dummyChannelName, i),
				profile.ID,
				guild.ID,
				s.channelRepo,
			)
			s.NoError(err)
			createdChannels = append(createdChannels, *channel)
		}

		payload := endpoint.FetchGuildChannelsRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/channels?page=2", guild.ID)
		req, err := http.NewRequest(http.MethodGet, s.Route(rawRoute), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusOK, res.StatusCode)

		var data endpoint.FetchAvailableGuildsResponse
		err = parsePayload[endpoint.FetchAvailableGuildsResponse](&data, res.Body)
		s.NoError(err)

		s.Equal("guild channels fetched", data.Message)
		s.Equal(totalChannels-core.ItemsPerPage(), len(data.Data))
		s.Equal(http.StatusOK, data.Code)

		for idx, g := range data.Data {
			realIdx := core.ItemsPerPage() + idx

			s.Equal(createdChannels[realIdx].ID, g.ID)
			s.Equal(createdChannels[realIdx].Name, g.Name)
			s.Equal(createdChannels[realIdx].GuildID, guild.ID)
		}
	})

	s.Run("it should be able fetch guild channels without passing page as param", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy guild name 1", s.guildRepo)
		s.NoError(err)
		_, err = createSampleMember(profile.ID, user.ID, guild.ID, s.guildMemberRepo)
		s.NoError(err)

		totalChannels := 4
		createdChannels := []entity.Channel{}
		for i := 0; i < totalChannels; i++ {
			channel, err := createSampleChannel(
				fmt.Sprintf("%s-%d", dummyChannelName, i),
				profile.ID,
				guild.ID,
				s.channelRepo,
			)
			s.NoError(err)
			createdChannels = append(createdChannels, *channel)
		}

		payload := endpoint.FetchGuildChannelsRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/channels", guild.ID)
		req, err := http.NewRequest(http.MethodGet, s.Route(rawRoute), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusOK, res.StatusCode)

		var data endpoint.FetchGuildChannelsResponse
		err = parsePayload[endpoint.FetchGuildChannelsResponse](&data, res.Body)
		s.NoError(err)

		s.Equal("guild channels fetched", data.Message)
		s.Equal(totalChannels, len(data.Data))
		s.Equal(http.StatusOK, data.Code)

		for idx, g := range data.Data {
			s.Equal(createdChannels[idx].ID, g.ID)
			s.Equal(createdChannels[idx].Name, g.Name)
			s.Equal(guild.ID, g.GuildID)
		}
	})

	s.Run("it should be not able to fetch guild channels with an invalid guild_id", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		payload := endpoint.FetchGuildChannelsRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, s.Route("/v1/guilds/-2/channels"), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusNotFound, res.StatusCode)

		var data endpoint.FetchGuildChannelsResponse
		err = parsePayload[endpoint.FetchGuildChannelsResponse](&data, res.Body)
		s.NoError(err)

		s.Equal(validation.NewNotFoundErr("guild").Error(), data.Message)
		s.Equal(http.StatusNotFound, data.Code)
	})

	s.Run("it should be able fetch guild channels with the max capacity", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy guild name 1", s.guildRepo)
		s.NoError(err)
		_, err = createSampleMember(profile.ID, user.ID, guild.ID, s.guildMemberRepo)
		s.NoError(err)

		totalChannels := core.ItemsPerPage() + 1
		createdChannels := []entity.Channel{}
		for i := 0; i < totalChannels; i++ {
			channel, err := createSampleChannel(
				fmt.Sprintf("%s-%d", dummyChannelName, i),
				profile.ID,
				guild.ID,
				s.channelRepo,
			)
			s.NoError(err)
			createdChannels = append(createdChannels, *channel)
		}

		payload := endpoint.FetchGuildChannelsRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/channels", guild.ID)
		req, err := http.NewRequest(http.MethodGet, s.Route(rawRoute), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusOK, res.StatusCode)

		var data endpoint.FetchGuildChannelsResponse
		err = parsePayload[endpoint.FetchGuildChannelsResponse](&data, res.Body)
		s.NoError(err)

		s.Equal("guild channels fetched", data.Message)
		s.Equal(core.ItemsPerPage(), len(data.Data))
		s.Equal(http.StatusOK, data.Code)

		for idx, g := range data.Data {
			s.Equal(createdChannels[idx].ID, g.ID)
			s.Equal(createdChannels[idx].Name, g.Name)
			s.Equal(guild.ID, g.GuildID)
		}
	})

	s.Run("it should be not able fetch guild channels by when passing an invalid page type", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy guild name 1", s.guildRepo)
		s.NoError(err)

		payload := endpoint.FetchGuildChannelsRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/channels?page=test", guild.ID)
		req, err := http.NewRequest(http.MethodGet, s.Route(rawRoute), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusBadRequest, res.StatusCode)

		var data endpoint.Response
		err = parsePayload[endpoint.Response](&data, res.Body)
		s.NoError(err)

		s.Equal("`test` is not a valid page", data.Message)
		s.Equal(http.StatusBadRequest, data.Code)
	})

	s.Run("it should be not able to fetch guild channels if profile dont exists", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy guild name", s.guildRepo)
		s.NoError(err)

		payload := endpoint.FetchGuildChannelsRequest{
			ProfileID: -2,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/channels", guild.ID)
		req, err := http.NewRequest(http.MethodGet, s.Route(rawRoute), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusNotFound, res.StatusCode)

		var data endpoint.Response
		err = parsePayload[endpoint.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(validation.NewNotFoundErr("user_profile").Error(), data.Message)
		s.Equal(http.StatusNotFound, data.Code)
	})

	s.Run("it should be not able to fetch guild channels if user dont exists", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)
		guild, err := createSampleGuild(profile.ID, "dummy guild name", s.guildRepo)
		s.NoError(err)

		payload := endpoint.FetchGuildChannelsRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/channels", guild.ID)
		req, err := http.NewRequest(http.MethodGet, s.Route(rawRoute), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, -2)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusNotFound, res.StatusCode)

		var data endpoint.Response
		err = parsePayload[endpoint.Response](&data, res.Body)
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

		payload := endpoint.FetchGuildChannelsRequest{
			ProfileID: profile.ID,
		}
		body, err := json.Marshal(payload)
		s.NoError(err)

		client := &http.Client{}
		rawRoute := fmt.Sprintf("/v1/guilds/%d/channels", guild.ID)
		req, err := http.NewRequest(http.MethodGet, s.Route(rawRoute), writeBody(body))
		s.NoError(err)
		err = authenticate(s.token, req, user.ID)
		s.NoError(err)
		res, err := client.Do(req)
		s.NoError(err)

		defer res.Body.Close()

		s.Equal(http.StatusUnauthorized, res.StatusCode)

		var data endpoint.Response
		err = parsePayload[endpoint.Response](&data, res.Body)
		s.NoError(err)

		s.Equal(validation.NewUnauthorizedErr().Error(), data.Message)
		s.Equal(http.StatusUnauthorized, data.Code)
	})
}
