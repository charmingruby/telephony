package integration

import (
	guildEntity "github.com/charmingruby/telephony/internal/domain/guild/entity"
	userEntity "github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/infra/database"
	"github.com/charmingruby/telephony/internal/infra/security/cryptography"
)

func createSampleUser(
	email string,
	userRepo *database.PostgresUserRepository,
) (*userEntity.User, error) {
	passwordHash, err := cryptography.NewCryptography().GenerateHash("password123")
	if err != nil {
		return nil, err
	}

	user, err := userEntity.NewUser(
		"dummy name",
		"dummy lastname",
		email,
		"password",
	)
	user.PasswordHash = passwordHash
	if err != nil {
		return nil, err
	}

	id, err := userRepo.Store(user)
	if err != nil {
		return nil, err
	}

	user.ID = id

	return user, nil
}

func createSampleUserProfile(
	userID int,
	displayName string,
	profileRepo *database.PostgresUserProfileRepository,
) (*userEntity.UserProfile, error) {
	profile, err := userEntity.NewUserProfile(
		displayName,
		"dummy biography",
		userID,
	)
	if err != nil {
		return nil, err
	}

	id, err := profileRepo.Store(profile)
	if err != nil {
		return nil, err
	}

	profile.ID = id

	return profile, nil
}

func createSampleGuild(
	profileID int,
	name string,
	guildRepo *database.PostgresGuildRepository,
) (*guildEntity.Guild, error) {
	guild, err := guildEntity.NewGuild(
		name,
		"dummy biography",
		profileID,
	)
	if err != nil {
		return nil, err
	}

	id, err := guildRepo.Store(guild)
	if err != nil {
		return nil, err
	}

	guild.ID = id

	return guild, nil
}
