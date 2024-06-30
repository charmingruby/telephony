package usecase

import (
	"fmt"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	guildEntity "github.com/charmingruby/telephony/internal/domain/guild/entity"
	userEntity "github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_FetchGuildChannels() {
	user, err := userEntity.NewUser("dummy name", "dummy last name", "dummy@email.com", "password123")
	s.NoError(err)

	profile, err := userEntity.NewUserProfile("dummy name", "dummy biography", user.ID)
	s.NoError(err)

	guild, err := guildEntity.NewGuild(
		"dummy name",
		"dummy description",
		profile.ID,
	)
	s.NoError(err)

	dummyChannelName := "dummy name"

	s.Run("it should be able to fetch guild channels", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		for i := 0; i < core.ItemsPerPage()+1; i++ {
			ch, err := guildEntity.NewChannel(
				fmt.Sprintf("%s-%d", dummyChannelName, i+1),
				profile.ID,
				guild.ID,
			)
			s.NoError(err)

			_, err = s.channelRepo.Store(ch)
			s.NoError(err)
		}
		s.Equal(core.ItemsPerPage()+1, len(s.channelRepo.Items))

		dto := dto.FetchGuildChannelsDTO{
			UserID:    user.ID,
			ProfileID: profile.ID,
			GuildID:   guild.ID,
			Pagination: core.PaginationParams{
				Page: 2,
			},
		}

		channels, err := s.guildService.FetchGuildChannels(dto)

		s.NoError(err)
		s.Equal(1, len(channels))
		s.Equal(fmt.Sprintf("%s-%d", dummyChannelName, 51), channels[0].Name)
		s.Equal(guild.ID, channels[0].GuildID)
	})

	s.Run("it should be able to fetch guild channels with max capacity", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		for i := 0; i < core.ItemsPerPage()*2; i++ {
			if i%2 == 0 {
				ch, err := guildEntity.NewChannel(
					fmt.Sprintf("%s-%d", dummyChannelName, i+1),
					profile.ID,
					-2,
				)
				s.NoError(err)

				_, err = s.channelRepo.Store(ch)
				s.NoError(err)
			} else {
				ch, err := guildEntity.NewChannel(
					fmt.Sprintf("%s-%d", dummyChannelName, i+1),
					profile.ID,
					guild.ID,
				)
				s.NoError(err)

				_, err = s.channelRepo.Store(ch)
				s.NoError(err)
			}
		}
		s.Equal(core.ItemsPerPage()*2, len(s.channelRepo.Items))

		dto := dto.FetchGuildChannelsDTO{
			UserID:    user.ID,
			ProfileID: profile.ID,
			GuildID:   guild.ID,
			Pagination: core.PaginationParams{
				Page: 1,
			},
		}

		channels, err := s.guildService.FetchGuildChannels(dto)

		s.NoError(err)
		s.Equal(core.ItemsPerPage(), len(channels))

		for _, ch := range channels {
			s.Equal(guild.ID, ch.GuildID)
		}
	})

	s.Run("it should be not able to fetch guilds channels if guild dont exists", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		dto := dto.FetchGuildChannelsDTO{
			UserID:    user.ID,
			ProfileID: profile.ID,
			GuildID:   -2,
			Pagination: core.PaginationParams{
				Page: 1,
			},
		}

		channels, err := s.guildService.FetchGuildChannels(dto)

		s.Error(err)
		s.Nil(channels)
		s.Equal(validation.NewNotFoundErr("guild").Error(), err.Error())
	})

	s.Run("it should be not able to create a guild if profile do not exists", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		dto := dto.FetchGuildChannelsDTO{
			ProfileID: -2,
			UserID:    user.ID,
			GuildID:   guild.ID,
			Pagination: core.PaginationParams{
				Page: 1,
			},
		}

		channels, err := s.guildService.FetchGuildChannels(dto)

		s.Error(err)
		s.Nil(channels)
		s.Equal(validation.NewNotFoundErr("user_profile").Error(), err.Error())
	})

	s.Run("it should be not able to create a guild if user do not exists", func() {
		_, err := s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		dto := dto.FetchGuildChannelsDTO{
			ProfileID: profile.ID,
			UserID:    -2,
			GuildID:   guild.ID,
			Pagination: core.PaginationParams{
				Page: 1,
			},
		}

		channels, err := s.guildService.FetchGuildChannels(dto)

		s.Error(err)
		s.Nil(channels)
		s.Equal(validation.NewNotFoundErr("user").Error(), err.Error())
	})

	s.Run("it should be not able to create a guild if its not the profile owner", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		otherUser := *user
		otherUser.ID = -2
		_, err = s.userRepo.Store(&otherUser)
		s.NoError(err)
		s.Equal(-2, s.userRepo.Items[1].ID)

		otherUserProfile := *profile
		otherUserProfile.ID = -2
		otherUserProfile.UserID = otherUser.ID
		_, err = s.profileRepo.Store(&otherUserProfile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		dto := dto.FetchGuildChannelsDTO{
			ProfileID: otherUserProfile.ID,
			UserID:    user.ID,
			GuildID:   guild.ID,
			Pagination: core.PaginationParams{
				Page: 1,
			},
		}

		channels, err := s.guildService.FetchGuildChannels(dto)

		s.Error(err)
		s.Nil(channels)
		s.Equal(validation.NewUnauthorizedErr().Error(), err.Error())
	})
}
