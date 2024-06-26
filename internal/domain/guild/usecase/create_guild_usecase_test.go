package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	guildEntity "github.com/charmingruby/telephony/internal/domain/guild/entity"
	userEntity "github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_CreateGuild() {
	user, err := userEntity.NewUser("dummy name", "dummy last name", "dummy@email.com", "password123")
	s.NoError(err)

	profile, err := userEntity.NewUserProfile("dummy name", "dummy biography", user.ID)
	s.NoError(err)

	dummyGuildName := "dummy name"
	dummyGuildDescription := "dummy description"
	dummyTags := []string{"Development"}

	guild, err := guildEntity.NewGuild(dummyGuildName, dummyGuildDescription, dummyTags, profile.ID)
	s.NoError(err)

	s.Run("it should be able to create a guild", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		dto := dto.CreateGuildDTO{
			Name:        dummyGuildName,
			Description: dummyGuildDescription,
			Tags:        dummyTags,
			ProfileID:   profile.ID,
			UserID:      user.ID,
		}

		err = s.guildService.CreateGuild(dto)

		s.NoError(err)
		s.Equal(dummyGuildName, s.guildRepo.Items[0].Name)
	})

	s.Run("it should be not able to create a guild if profile do not exists", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		dto := dto.CreateGuildDTO{
			Name:        dummyGuildName,
			Description: dummyGuildDescription,
			Tags:        dummyTags,
			ProfileID:   -2,
			UserID:      user.ID,
		}

		err = s.guildService.CreateGuild(dto)

		s.Error(err)
		s.Equal(validation.NewNotFoundErr("user_profile").Error(), err.Error())
	})

	s.Run("it should be not able to create a guild if user do not exists", func() {
		_, err := s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		dto := dto.CreateGuildDTO{
			Name:        dummyGuildName,
			Description: dummyGuildDescription,
			Tags:        dummyTags,
			ProfileID:   profile.ID,
			UserID:      -2,
		}

		err = s.guildService.CreateGuild(dto)

		s.Error(err)
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

		dto := dto.CreateGuildDTO{
			Name:        dummyGuildName,
			Description: dummyGuildDescription,
			Tags:        dummyTags,
			ProfileID:   otherUserProfile.ID,
			UserID:      user.ID,
		}

		err = s.guildService.CreateGuild(dto)

		s.Error(err)
		s.Equal(validation.NewUnauthorizedErr().Error(), err.Error())
	})

	s.Run("it should be not able to create a guild if name is already taken", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
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
			UserID:      user.ID,
		}

		err = s.guildService.CreateGuild(dto)

		s.Error(err)
		s.Equal(validation.NewConflictErr("guild", "name").Error(), err.Error())
	})

	s.Run("it should be not able to create a guild with validation error", func() {
		_, err = s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err := s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		dto := dto.CreateGuildDTO{
			Name:        "",
			Description: dummyGuildDescription,
			Tags:        dummyTags,
			ProfileID:   profile.ID,
			UserID:      user.ID,
		}

		err = s.guildService.CreateGuild(dto)

		s.Error(err)
		s.Equal(validation.NewValidationErr(validation.ErrRequired("name")).Error(), err.Error())
	})
}
