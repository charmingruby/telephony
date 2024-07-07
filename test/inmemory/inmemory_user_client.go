package inmemory

import (
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func NewInMemoryUserProfileClient(
	profileRepo *InMemoryUserProfileRepository,
	userRepo *InMemoryUserRepository,
) *InMemoryUserProfileClient {
	return &InMemoryUserProfileClient{
		Items:       []entity.UserProfile{},
		ProfileRepo: profileRepo,
		UserRepo:    userRepo,
	}
}

type InMemoryUserProfileClient struct {
	Items       []entity.UserProfile
	ProfileRepo *InMemoryUserProfileRepository
	UserRepo    *InMemoryUserRepository
}

func (c *InMemoryUserProfileClient) ProfileExists(id int) bool {
	_, err := c.ProfileRepo.FindByID(id)
	return err == nil
}

func (c *InMemoryUserProfileClient) UserExists(id int) bool {
	_, err := c.UserRepo.FindByID(id)
	return err == nil
}

func (c *InMemoryUserProfileClient) GetDisplayName(profileID int) (string, error) {
	profile, err := c.ProfileRepo.FindByID(profileID)
	if err != nil {
		return "", validation.NewNotFoundErr("profile")
	}

	return profile.DisplayName, nil
}

func (c *InMemoryUserProfileClient) IsTheProfileOwner(userID, profileID int) bool {
	profile, err := c.ProfileRepo.FindByID(profileID)
	if err != nil {
		return false
	}

	return profile.UserID == userID
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
