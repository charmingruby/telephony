package inmemory

import "github.com/charmingruby/telephony/internal/domain/user/entity"

func NewInMemoryUserProfileClient(profileRepo *InMemoryUserProfileRepository) *InMemoryUserProfileClient {
	return &InMemoryUserProfileClient{
		Items:       []entity.UserProfile{},
		ProfileRepo: profileRepo,
	}
}

type InMemoryUserProfileClient struct {
	Items       []entity.UserProfile
	ProfileRepo *InMemoryUserProfileRepository
}

func (c *InMemoryUserProfileClient) ProfileExists(id int) bool {
	_, err := c.ProfileRepo.FindByID(id)
	return err == nil
}

func (c *InMemoryUserProfileClient) GuildJoin(id int) error {
	if err := c.ProfileRepo.UpdateGuildsQuantity(id, +1); err != nil {
		return err
	}

	return nil
}

func (c *InMemoryUserProfileClient) GuildLeave(id int) error {
	if err := c.ProfileRepo.UpdateGuildsQuantity(id, -1); err != nil {
		return err
	}

	return nil
}

func (c *InMemoryUserProfileClient) SendMessage(id int) error {
	if err := c.ProfileRepo.UpdateMessagesQuantity(id, +1); err != nil {
		return err
	}

	return nil
}
