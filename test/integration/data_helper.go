package integration

import (
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/infra/database"
	"github.com/charmingruby/telephony/internal/infra/security/cryptography"
)

func createSampleUser(
	userRepo *database.PostgresUserRepository,
	profileRepo *database.PostgresUserProfileRepository,
) (*entity.User, *entity.UserProfile, error) {
	passwordHash, err := cryptography.NewCryptography().GenerateHash("password123")
	if err != nil {
		return nil, nil, err
	}

	user, err := entity.NewUser(
		"dummy name",
		"dummy lastname",
		"dummy@email.com",
		"password",
	)
	user.PasswordHash = passwordHash
	if err != nil {
		return nil, nil, err
	}

	if err := userRepo.Store(user); err != nil {
		return nil, nil, err
	}

	userPersisted, err := userRepo.FindByEmail(user.Email)
	if err != nil {
		return nil, nil, err
	}
	user.ID = userPersisted.ID

	profile, err := entity.NewUserProfile(
		"dummy nick",
		"dummy biography",
		userPersisted.ID,
	)
	if err != nil {
		return nil, nil, err
	}

	if err := profileRepo.Store(profile); err != nil {
		return nil, nil, err
	}

	return user, profile, nil
}
