package usecase

import (
	"fmt"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	userEntity "github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_FetchGuildChannels() {
	profile, err := userEntity.NewUserProfile("dummy name", "dummy biography", 1)
	s.NoError(err)

	guild, err := entity.NewGuild(
		"dummy name",
		"dummy description",
		profile.ID,
	)
	s.NoError(err)

	dummyChannelName := "dummy name"

	s.Run("it should be able to fetch guild channels", func() {
		_, err := s.guildRepo.Store(guild)
		s.NoError(err)

		for i := 0; i < core.ItemsPerPage()+1; i++ {
			ch, err := entity.NewChannel(
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
			GuildID: guild.ID,
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
		_, err := s.guildRepo.Store(guild)
		s.NoError(err)

		for i := 0; i < core.ItemsPerPage()*2; i++ {
			if i%2 == 0 {
				ch, err := entity.NewChannel(
					fmt.Sprintf("%s-%d", dummyChannelName, i+1),
					profile.ID,
					-2,
				)
				s.NoError(err)

				_, err = s.channelRepo.Store(ch)
				s.NoError(err)
			} else {
				ch, err := entity.NewChannel(
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
			GuildID: guild.ID,
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
		dto := dto.FetchGuildChannelsDTO{
			GuildID: -2,
			Pagination: core.PaginationParams{
				Page: 1,
			},
		}

		channels, err := s.guildService.FetchGuildChannels(dto)

		s.Error(err)
		s.Nil(channels)
		s.Equal(validation.NewNotFoundErr("guild").Error(), err.Error())
	})
}
