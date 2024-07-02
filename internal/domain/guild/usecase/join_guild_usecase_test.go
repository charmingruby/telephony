package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	guildEntity "github.com/charmingruby/telephony/internal/domain/guild/entity"
	userEntity "github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_JoinGuild() {
	user, err := userEntity.NewUser("dummy name", "dummy last name", "dummy@email.com", "password123")
	s.NoError(err)

	profile, err := userEntity.NewUserProfile("dummy name", "dummy biography", user.ID)
	s.NoError(err)

	guild, err := guildEntity.NewGuild("dummy name", "dummy description", profile.ID)
	s.NoError(err)

	member, err := guildEntity.NewGuildMember(profile.ID, user.ID, guild.ID)
	s.NoError(err)

	s.Run("it should be able to join a guild successfully", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		s.Equal(0, len(s.memberRepo.Items))

		dto := dto.JoinGuildDTO{
			ProfileID: profile.ID,
			UserID:    user.ID,
			GuildID:   guild.ID,
		}

		err = s.guildService.JoinGuild(dto)

		s.NoError(err)
		s.Equal(1, len(s.memberRepo.Items))
		s.Equal(dto.UserID, s.memberRepo.Items[0].UserID)
		s.Equal(dto.ProfileID, s.memberRepo.Items[0].ProfileID)
		s.Equal(dto.GuildID, s.memberRepo.Items[0].GuildID)
	})

	s.Run("it should be not able to join a guild if profile do not exists", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		dto := dto.JoinGuildDTO{
			ProfileID: -2,
			UserID:    user.ID,
			GuildID:   guild.ID,
		}

		err = s.guildService.JoinGuild(dto)

		s.Error(err)
		s.Equal(validation.NewNotFoundErr("user_profile").Error(), err.Error())
	})

	s.Run("it should be not able to create a channel if user do not exists", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		dto := dto.JoinGuildDTO{
			ProfileID: profile.ID,
			UserID:    -2,
			GuildID:   guild.ID,
		}

		err = s.guildService.JoinGuild(dto)

		s.Error(err)
		s.Equal(validation.NewNotFoundErr("user").Error(), err.Error())
	})

	s.Run("it should be not able to join a guild if its not the profile owner", func() {
		randomUser := *user
		randomUser.ID = -2
		_, err := s.userRepo.Store(&randomUser)
		s.NoError(err)
		_, err = s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(2, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		s.Equal(0, len(s.memberRepo.Items))

		dto := dto.JoinGuildDTO{
			ProfileID: profile.ID,
			UserID:    randomUser.ID,
			GuildID:   guild.ID,
		}

		err = s.guildService.JoinGuild(dto)

		s.Error(err)
		s.Equal(validation.NewUnauthorizedErr().Error(), err.Error())
	})

	s.Run("it should be not able to join a guild if is already a member", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		_, err = s.memberRepo.Store(member)
		s.NoError(err)
		s.Equal(1, len(s.memberRepo.Items))

		dto := dto.JoinGuildDTO{
			ProfileID: profile.ID,
			UserID:    user.ID,
			GuildID:   guild.ID,
		}

		err = s.guildService.JoinGuild(dto)

		s.Error(err)
		s.Equal(validation.NewBadRequestErr("user profile is already a member").Error(), err.Error())
	})
}
