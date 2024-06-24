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

	if err := userRepo.Store(user); err != nil {
		return nil, err
	}

	userPersisted, err := userRepo.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	user.ID = userPersisted.ID

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

	if err := profileRepo.Store(profile); err != nil {
		return nil, err
	}

	return profile, nil
}
