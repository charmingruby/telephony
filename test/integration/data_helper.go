package integration

import (
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/infra/database"
	"github.com/charmingruby/telephony/internal/infra/security/cryptography"
)

func createSampleUser(
	email string,
	userRepo *database.PostgresUserRepository,
) (*entity.User, error) {
	passwordHash, err := cryptography.NewCryptography().GenerateHash("password123")
	if err != nil {
		return nil, err
	}

	user, err := entity.NewUser(
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
) (*entity.UserProfile, error) {
	profile, err := entity.NewUserProfile(
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
