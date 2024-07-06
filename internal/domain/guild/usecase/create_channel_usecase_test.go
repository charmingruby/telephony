package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/guild/dto"
	guildEntity "github.com/charmingruby/telephony/internal/domain/guild/entity"
	userEntity "github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *Suite) Test_CreateChannel() {
	user, err := userEntity.NewUser("dummy name", "dummy last name", "dummy@email.com", "password123")
	s.NoError(err)

	profile, err := userEntity.NewUserProfile("dummy name", "dummy biography", user.ID)
	s.NoError(err)

	guild, err := guildEntity.NewGuild("dummy name", "dummy description", profile.ID)
	s.NoError(err)

	dummyChannelName := "dummy channel"
	channel, err := guildEntity.NewChannel(dummyChannelName, guild.ID, profile.ID)
	s.NoError(err)

	member, err := guildEntity.NewGuildMember(profile.ID, user.ID, guild.ID)
	s.NoError(err)

	s.Run("it should be able to create a channel", func() {
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

		dto := dto.CreateChannelDTO{
			Name:      dummyChannelName,
			GuildID:   guild.ID,
			ProfileID: profile.ID,
			UserID:    user.ID,
		}

		channelID, err := s.guildService.CreateChannel(dto)

		s.NoError(err)
		s.Equal(s.channelRepo.Items[0].ID, channelID)
		s.Equal(dummyChannelName, s.channelRepo.Items[0].Name)
	})

	s.Run("it should be not able to create a channel with conflicting name at the same guild", func() {
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

		_, err = s.memberRepo.Store(member)
		s.NoError(err)
		s.Equal(1, len(s.memberRepo.Items))

		dto := dto.CreateChannelDTO{
			Name:      dummyChannelName,
			GuildID:   guild.ID,
			ProfileID: profile.ID,
			UserID:    user.ID,
		}

		channelID, err := s.guildService.CreateChannel(dto)

		s.Error(err)
		s.Equal(-1, channelID)
		s.Equal(validation.NewConflictErr("channel", "name").Error(), err.Error())
		s.Equal(1, len(s.channelRepo.Items))
	})

	s.Run("it should be not able to create a channel if guild dont exists", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		s.Equal(0, len(s.guildRepo.Items))

		_, err = s.channelRepo.Store(channel)
		s.NoError(err)
		s.Equal(1, len(s.channelRepo.Items))

		dto := dto.CreateChannelDTO{
			Name:      dummyChannelName,
			GuildID:   -2,
			ProfileID: profile.ID,
			UserID:    user.ID,
		}

		channelID, err := s.guildService.CreateChannel(dto)

		s.Error(err)
		s.Equal(-1, channelID)
		s.Equal(validation.NewNotFoundErr("guild").Error(), err.Error())
		s.Equal(1, len(s.channelRepo.Items))

	})

	s.Run("it should be not able to create a channel if profile do not exists", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		dto := dto.CreateChannelDTO{
			Name:      dummyChannelName,
			GuildID:   guild.ID,
			ProfileID: -2,
			UserID:    user.ID,
		}

		channelID, err := s.guildService.CreateChannel(dto)

		s.Error(err)
		s.Equal(-1, channelID)
		s.Equal(validation.NewNotFoundErr("user_profile").Error(), err.Error())
	})

	s.Run("it should be not able to create a channel if user do not exists", func() {
		_, err := s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		dto := dto.CreateChannelDTO{
			Name:      dummyChannelName,
			GuildID:   guild.ID,
			ProfileID: -2,
			UserID:    user.ID,
		}

		channelID, err := s.guildService.CreateChannel(dto)

		s.Error(err)
		s.Equal(-1, channelID)
		s.Equal(validation.NewNotFoundErr("user").Error(), err.Error())
	})

	s.Run("it should be not able to create a channel if its not the profile owner", func() {
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

		dto := dto.CreateChannelDTO{
			Name:      dummyChannelName,
			GuildID:   guild.ID,
			ProfileID: otherUserProfile.ID,
			UserID:    user.ID,
		}

		channelID, err := s.guildService.CreateChannel(dto)

		s.Error(err)
		s.Equal(-1, channelID)
		s.Equal(validation.NewUnauthorizedErr().Error(), err.Error())
	})

	s.Run("it should be not able to create a channel with validation error", func() {
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

		dto := dto.CreateChannelDTO{
			Name:      "",
			GuildID:   guild.ID,
			ProfileID: profile.ID,
			UserID:    user.ID,
		}

		channelID, err := s.guildService.CreateChannel(dto)

		s.Error(err)
		s.Equal(-1, channelID)
		s.Equal(validation.NewValidationErr(validation.ErrRequired("name")).Error(), err.Error())
	})

	s.Run("it should be not able to create a channel if its not the guild owner", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		randomProfile := *profile
		randomProfile.ID = -2
		_, err = s.profileRepo.Store(&randomProfile)
		s.NoError(err)
		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(2, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		randomMember := *member
		randomMember.ID = -2
		randomMember.ProfileID = randomProfile.ID
		_, err = s.memberRepo.Store(&randomMember)
		s.NoError(err)
		_, err = s.memberRepo.Store(member)
		s.NoError(err)
		s.Equal(2, len(s.memberRepo.Items))

		dto := dto.CreateChannelDTO{
			Name:      dummyChannelName,
			GuildID:   guild.ID,
			ProfileID: randomProfile.ID,
			UserID:    user.ID,
		}

		channelID, err := s.guildService.CreateChannel(dto)

		s.Error(err)
		s.Equal(-1, channelID)
		s.Equal(validation.NewUnauthorizedErr().Error(), err.Error())
	})

	s.Run("it should be not able to create a channel if its not a guild member", func() {
		_, err := s.userRepo.Store(user)
		s.NoError(err)
		s.Equal(1, len(s.userRepo.Items))

		_, err = s.profileRepo.Store(profile)
		s.NoError(err)
		s.Equal(1, len(s.profileRepo.Items))

		_, err = s.guildRepo.Store(guild)
		s.NoError(err)
		s.Equal(1, len(s.guildRepo.Items))

		dto := dto.CreateChannelDTO{
			Name:      dummyChannelName,
			GuildID:   guild.ID,
			ProfileID: profile.ID,
			UserID:    user.ID,
		}

		channelID, err := s.guildService.CreateChannel(dto)

		s.Error(err)
		s.Equal(-1, channelID)
		s.Equal(validation.NewUnauthorizedErr().Error(), err.Error())
	})
}
