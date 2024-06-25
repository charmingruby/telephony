package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	guildEntity "github.com/charmingruby/telephony/internal/domain/guild/entity"
	userEntity "github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_CreateGuild() {
	profile, err := userEntity.NewUserProfile("dummy name", "dummy biography", 1)
	s.NoError(err)

	dummyGuildName := "dummy name"
	dummyGuildDescription := "dummy description"
	dummyTags := []string{"Development"}

	guild, err := guildEntity.NewGuild(dummyGuildName, dummyGuildDescription, dummyTags, profile.ID)
	s.NoError(err)

	s.Run("it should be able to create a guild", func() {
		_, err := s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		dto := dto.CreateGuildDTO{
			Name:        dummyGuildName,
			Description: dummyGuildDescription,
			Tags:        dummyTags,
			ProfileID:   profile.ID,
		}

		err = s.guildService.CreateGuild(dto)

		s.NoError(err)
		s.Equal(dummyGuildName, s.guildRepo.Items[0].Name)
	})

	s.Run("it should be not able to create a guild if owner don't exists", func() {
		s.Equal(0, len(s.profileRepo.Items))

		dto := dto.CreateGuildDTO{
			Name:        dummyGuildName,
			Description: dummyGuildDescription,
			Tags:        dummyTags,
			ProfileID:   1,
		}

		err = s.guildService.CreateGuild(dto)

		s.Error(err)
		s.Equal(validation.NewNotFoundErr("user_profile").Error(), err.Error())
	})

	s.Run("it should be not able to create a guild if name is already taken", func() {
		_, err := s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		dto := dto.CreateGuildDTO{
			Name:        guild.Name,
			Description: dummyGuildDescription,
			Tags:        dummyTags,
			ProfileID:   profile.ID,
		}

		err = s.guildService.CreateGuild(dto)

		s.Error(err)
		s.Equal(validation.NewConflictErr("guild", "name").Error(), err.Error())
	})

	s.Run("it should be not able to create a guild with validation error", func() {
		_, err := s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		dto := dto.CreateGuildDTO{
			Name:        "",
			Description: dummyGuildDescription,
			Tags:        dummyTags,
			ProfileID:   profile.ID,
		}

		err = s.guildService.CreateGuild(dto)

		s.Error(err)
		s.Equal(validation.NewValidationErr(validation.ErrRequired("name")).Error(), err.Error())
	})
}
