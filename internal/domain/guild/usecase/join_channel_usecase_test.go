package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	guildEntity "github.com/charmingruby/telephony/internal/domain/guild/entity"
	userEntity "github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_JoinChannel() {
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

	member, err := guildEntity.NewGuildMember(profile.ID, user.ID, guild.ID)
	s.NoError(err)

	channel, err := guildEntity.NewChannel("dummy channel", guild.ID, profile.ID)
	s.NoError(err)

	s.Run("it should be able to join channel", func() {
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

		_, err = s.channelRepo.Store(channel)
		s.NoError(err)
		s.Equal(1, len(s.channelRepo.Items))

		dto := dto.JoinChannelDTO{
			ChannelID: channel.ID,
			UserID:    user.ID,
			ProfileID: profile.ID,
			GuildID:   guild.ID,
		}

		res, err := s.guildService.JoinChannel(dto)

		s.NoError(err)
		s.Equal(profile.DisplayName, res.DisplayName)
	})

	s.Run("it should be not able to join channel if guild dont exists", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.memberRepo.Store(member)
		s.NoError(err)
		s.Equal(1, len(s.memberRepo.Items))

		_, err = s.channelRepo.Store(channel)
		s.NoError(err)
		s.Equal(1, len(s.channelRepo.Items))

		dto := dto.JoinChannelDTO{
			ChannelID: channel.ID,
			UserID:    user.ID,
			ProfileID: profile.ID,
			GuildID:   guild.ID,
		}

		res, err := s.guildService.JoinChannel(dto)

		s.Error(err)
		s.Nil(res)
		s.Equal(validation.NewNotFoundErr("guild").Error(), err.Error())
	})

	s.Run("it should be not able to join channel if profile do not exists", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.channelRepo.Store(channel)
		s.NoError(err)
		s.Equal(1, len(s.channelRepo.Items))

		dto := dto.JoinChannelDTO{
			ChannelID: channel.ID,
			UserID:    user.ID,
			ProfileID: -2,
			GuildID:   guild.ID,
		}

		res, err := s.guildService.JoinChannel(dto)

		s.Error(err)
		s.Nil(res)
		s.Equal(validation.NewNotFoundErr("user_profile").Error(), err.Error())
	})

	s.Run("it should be not able to join channel if user do not exists", func() {
		_, err := s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.channelRepo.Store(channel)
		s.NoError(err)
		s.Equal(1, len(s.channelRepo.Items))

		dto := dto.JoinChannelDTO{
			ChannelID: channel.ID,
			UserID:    -2,
			ProfileID: profile.ID,
			GuildID:   guild.ID,
		}

		res, err := s.guildService.JoinChannel(dto)

		s.Error(err)
		s.Nil(res)
		s.Equal(validation.NewNotFoundErr("user").Error(), err.Error())
	})

	s.Run("it should be not able to join channel if channel do not exists", func() {
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

		dto := dto.JoinChannelDTO{
			ChannelID: -2,
			UserID:    user.ID,
			ProfileID: profile.ID,
			GuildID:   guild.ID,
		}

		res, err := s.guildService.JoinChannel(dto)

		s.Error(err)
		s.Nil(res)
		s.Equal(validation.NewNotFoundErr("channel").Error(), err.Error())
	})

	s.Run("it should be not to join channel if its not a member", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		_, err = s.channelRepo.Store(channel)
		s.NoError(err)
		s.Equal(1, len(s.channelRepo.Items))

		dto := dto.JoinChannelDTO{
			ChannelID: channel.ID,
			UserID:    user.ID,
			ProfileID: profile.ID,
			GuildID:   guild.ID,
		}

		res, err := s.guildService.JoinChannel(dto)

		s.Error(err)
		s.Nil(res)
		s.Equal(validation.NewUnauthorizedErr().Error(), err.Error())
	})

	s.Run("it should be not able to join channel if its not the profile owner", func() {
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

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		dto := dto.JoinChannelDTO{
			ChannelID: channel.ID,
			UserID:    user.ID,
			ProfileID: otherUserProfile.ID,
			GuildID:   guild.ID,
		}

		res, err := s.guildService.JoinChannel(dto)

		s.Error(err)
		s.Nil(res)
		s.Equal(validation.NewUnauthorizedErr().Error(), err.Error())
	})
}
