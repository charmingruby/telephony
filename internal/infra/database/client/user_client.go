package client

import (
	"github.com/charmingruby/telephony/internal/infra/database"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewUserClient(
	profileRepo *database.PostgresUserProfileRepository,
	userRepo *database.PostgresUserRepository,
) *UserClient {
	return &UserClient{
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}

type UserClient struct {
	profileRepo *database.PostgresUserProfileRepository
	userRepo    *database.PostgresUserRepository
}

func (c *UserClient) UserExists(id int) bool {
	_, err := c.userRepo.FindByID(id)
	return err == nil
}

func (c *UserClient) GetDisplayName(profileID int) (string, error) {
	profile, err := c.profileRepo.FindByID(profileID)
	if err != nil {
		return "", validation.NewNotFoundErr("profile")
	}

	return profile.DisplayName, nil
}

func (c *UserClient) ProfileExists(id int) bool {
	_, err := c.profileRepo.FindByID(id)
	return err == nil
}

func (c *UserClient) IsTheProfileOwner(userID, profileID int) bool {
	profile, err := c.profileRepo.FindByID(profileID)
	if err != nil {
		return false
	}

	return profile.UserID == profileID
}

func (c *UserClient) GuildJoin(id int) error {
	if err := c.profileRepo.UpdateGuildsQuantity(id, +1); err != nil {
		return err
	}

	return nil
}

func (c *UserClient) GuildLeave(id int) error {
	if err := c.profileRepo.UpdateGuildsQuantity(id, -1); err != nil {
		return err
	}

	return nil
}

func (c *UserClient) SendMessage(id int) error {
	if err := c.profileRepo.UpdateMessagesQuantity(id, +1); err != nil {
		return err
	}

	return nil
}
