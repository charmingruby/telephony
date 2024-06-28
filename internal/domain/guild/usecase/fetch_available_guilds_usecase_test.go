package usecase

import (
	"fmt"
	"time"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	userEntity "github.com/charmingruby/telephony/internal/domain/user/entity"
)

func (s *Suite) Test_FetchAvailableGuilds() {
	profile, err := userEntity.NewUserProfile("dummy name", "dummy biography", 1)
	s.NoError(err)

	dummyName := "dummy name"
	dummyDescription := "dummy description"

	s.Run("it should be able to fetch paginated available guilds", func() {
		_, err := s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		for i := 0; i < 100; i++ {
			guild, _ := entity.NewGuild(
				fmt.Sprintf("%s%d", dummyName, i),
				dummyDescription,
				profile.ID,
			)
			guild.ID = i

			if i%2 == 0 {
				now := time.Now()
				guild.DeletedAt = &now
			}

			_, err = s.guildRepo.Store(guild)
			s.NoError(err)
		}

		s.Equal(100, len(s.guildRepo.Items))

		guilds, err := s.guildService.FetchAvailableGuilds(core.PaginationParams{Page: 1})
		s.NoError(err)
		s.Equal(25, len(guilds))

		for _, g := range guilds {
			isAvailable := g.DeletedAt == nil
			s.Equal(true, isAvailable)
		}
	})

	s.Run("it should be able to fetch available guilds even with the page is with empty guilds", func() {
		_, err := s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		for i := 0; i < 50; i++ {
			guild, _ := entity.NewGuild(
				fmt.Sprintf("%s%d", dummyName, i),
				dummyDescription,
				profile.ID,
			)
			guild.ID = i

			_, err = s.guildRepo.Store(guild)
			s.NoError(err)
		}

		s.Equal(50, len(s.guildRepo.Items))

		guilds, err := s.guildService.FetchAvailableGuilds(core.PaginationParams{Page: 2})
		s.NoError(err)
		s.Equal(0, len(guilds))
	})

	s.Run("it should be able to fetch available guilds overflowing", func() {
		_, err := s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		for i := 0; i < 53; i++ {
			guild, _ := entity.NewGuild(
				fmt.Sprintf("%s%d", dummyName, i),
				dummyDescription,
				profile.ID,
			)
			guild.ID = i

			_, err = s.guildRepo.Store(guild)
			s.NoError(err)
		}

		s.Equal(53, len(s.guildRepo.Items))

		guilds, err := s.guildService.FetchAvailableGuilds(core.PaginationParams{Page: 2})
		s.NoError(err)
		s.Equal(3, len(guilds))
	})
}
