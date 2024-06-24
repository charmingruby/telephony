package client

import "github.com/charmingruby/telephony/internal/infra/database"

func NewUserProfileClient(profileRepo *database.PostgresUserProfileRepository) *UserProfileClient {
	return &UserProfileClient{
		profileRepo: profileRepo,
	}
}

type UserProfileClient struct {
	profileRepo *database.PostgresUserProfileRepository
}

func (c *UserProfileClient) ProfileExists(id int) bool {
	_, err := c.profileRepo.FindByID(id)
	return err == nil
}

func (c *UserProfileClient) GuildJoin(id int) error {
	if err := c.profileRepo.UpdateGuildsQuantity(id, +1); err != nil {
		return err
	}

	return nil
}

func (c *UserProfileClient) GuildLeave(id int, quantityToDec int) error {
	if err := c.profileRepo.UpdateGuildsQuantity(id, -1); err != nil {
		return err
	}

	return nil
}

func (c *UserProfileClient) SendMessage(id int) error {
	if err := c.profileRepo.UpdateMessagesQuantity(id, +1); err != nil {
		return err
	}

	return nil
}
