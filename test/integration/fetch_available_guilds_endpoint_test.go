package integration

import (
	"fmt"
	"net/http"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
)

func (s *Suite) Test_FetchAvailableGuildsEndpoint() {
	dummyName := "dummy name"

	s.Run("it should be able fetch available guilds by page", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		totalGuilds := 4
		createdGuilds := []entity.Guild{}
		for i := 0; i < totalGuilds; i++ {
			guild, err := createSampleGuild(profile.ID, fmt.Sprintf("%s-%d", dummyName, i), s.guildRepo)
			s.NoError(err)
			createdGuilds = append(createdGuilds, *guild)
		}

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, s.Route("/v1/guilds?page=1"), nil)
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

		s.Equal("available guilds fetched", data.Message)
		s.Equal(totalGuilds, len(data.Data))
		s.Equal(http.StatusOK, data.Code)

		for idx, g := range data.Data {
			s.Equal(createdGuilds[idx].ID, g.ID)
			s.Equal(createdGuilds[idx].Name, g.Name)
		}
	})

	s.Run("it should be able fetch available guilds by page", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		totalGuilds := 4
		createdGuilds := []entity.Guild{}
		for i := 0; i < totalGuilds; i++ {
			if i%2 == 0 {
				guild, err := createSampleGuild(profile.ID, fmt.Sprintf("%s-%d", dummyName, i), s.guildRepo)
				s.NoError(err)
				createdGuilds = append(createdGuilds, *guild)
			}
		}

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, s.Route("/v1/guilds?page=1"), nil)
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

		s.Equal("available guilds fetched", data.Message)
		s.Equal(totalGuilds/2, len(data.Data))
		s.Equal(http.StatusOK, data.Code)

		for idx, g := range data.Data {
			s.Equal(createdGuilds[idx].ID, g.ID)
			s.Equal(createdGuilds[idx].Name, g.Name)
		}
	})

	s.Run("it should be able fetch available guilds even if is passing an empty page", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, s.Route("/v1/guilds?page=1"), nil)
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

		s.Equal("available guilds fetched", data.Message)
		s.Equal(0, len(data.Data))
		s.Equal(http.StatusOK, data.Code)
	})

	s.Run("it should be able fetch available guilds in a different page from the default", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		totalGuilds := core.ItemsPerPage() + 4
		createdGuilds := []entity.Guild{}
		for i := 0; i < totalGuilds; i++ {
			guild, err := createSampleGuild(profile.ID, fmt.Sprintf("%s-%d", dummyName, i), s.guildRepo)
			s.NoError(err)
			createdGuilds = append(createdGuilds, *guild)
		}

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, s.Route("/v1/guilds?page=2"), nil)
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

		s.Equal("available guilds fetched", data.Message)
		s.Equal(totalGuilds-core.ItemsPerPage(), len(data.Data))
		s.Equal(http.StatusOK, data.Code)

		for idx, g := range data.Data {
			realIdx := core.ItemsPerPage() + idx

			s.Equal(createdGuilds[realIdx].ID, g.ID)
			s.Equal(createdGuilds[realIdx].Name, g.Name)
		}
	})

	s.Run("it should be able fetch available guilds without passing page as param", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		totalGuilds := 4
		createdGuilds := []entity.Guild{}
		for i := 0; i < totalGuilds; i++ {
			guild, err := createSampleGuild(profile.ID, fmt.Sprintf("%s-%d", dummyName, i), s.guildRepo)
			s.NoError(err)
			createdGuilds = append(createdGuilds, *guild)
		}

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, s.Route("/v1/guilds"), nil)
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

		s.Equal("available guilds fetched", data.Message)
		s.Equal(totalGuilds, len(data.Data))
		s.Equal(http.StatusOK, data.Code)

		for idx, g := range data.Data {
			s.Equal(createdGuilds[idx].ID, g.ID)
			s.Equal(createdGuilds[idx].Name, g.Name)
		}
	})

	s.Run("it should be able fetch available guilds with the max capacity", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		totalGuilds := core.ItemsPerPage() + 1
		createdGuilds := []entity.Guild{}
		for i := 0; i < totalGuilds; i++ {
			guild, err := createSampleGuild(profile.ID, fmt.Sprintf("%s-%d", dummyName, i), s.guildRepo)
			s.NoError(err)
			createdGuilds = append(createdGuilds, *guild)
		}
		s.Equal(totalGuilds, len(createdGuilds))

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, s.Route("/v1/guilds?page=1"), nil)
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

		s.Equal("available guilds fetched", data.Message)
		s.Equal(core.ItemsPerPage(), len(data.Data))
		s.Equal(http.StatusOK, data.Code)

		for idx, g := range data.Data {
			s.Equal(createdGuilds[idx].ID, g.ID)
			s.Equal(createdGuilds[idx].Name, g.Name)
		}
	})

	s.Run("it should be not able fetch available guilds by when passing an invalid page type", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, s.Route("/v1/guilds?page=test"), nil)
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

	s.Run("it should be able fetch available guilds even if when is an zero or negative page value", func() {
		user, err := createSampleUser("dummy@email.com", s.userRepo)
		s.NoError(err)
		profile, err := createSampleUserProfile(user.ID, "dummy nick", s.profileRepo)
		s.NoError(err)

		totalGuilds := 4
		createdGuilds := []entity.Guild{}
		for i := 0; i < totalGuilds; i++ {
			guild, err := createSampleGuild(profile.ID, fmt.Sprintf("%s-%d", dummyName, i), s.guildRepo)
			s.NoError(err)
			createdGuilds = append(createdGuilds, *guild)
		}

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, s.Route("/v1/guilds?page=-2"), nil)
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

		s.Equal("available guilds fetched", data.Message)
		s.Equal(totalGuilds, len(data.Data))
		s.Equal(http.StatusOK, data.Code)

		for idx, g := range data.Data {
			s.Equal(createdGuilds[idx].ID, g.ID)
			s.Equal(createdGuilds[idx].Name, g.Name)
		}
	})
}
