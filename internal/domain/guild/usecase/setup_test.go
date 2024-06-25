package usecase

import (
	"testing"

	guildEntity "github.com/charmingruby/telephony/internal/domain/guild/entity"
	userEntity "github.com/charmingruby/telephony/internal/domain/user/entity"

	"github.com/charmingruby/telephony/test/inmemory"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	profileClient *inmemory.InMemoryUserProfileClient
	profileRepo   *inmemory.InMemoryUserProfileRepository
	guildRepo     *inmemory.InMemoryGuildRepository
	guildService  *GuildService
}

func (s *Suite) SetupSuite() {
	s.guildRepo = inmemory.NewInMemoryGuildRepository()
	s.profileRepo = inmemory.NewInMemoryUserProfileRepository()
	s.profileClient = inmemory.NewInMemoryUserProfileClient(s.profileRepo)
	s.guildService = NewGuildService(s.guildRepo, s.profileClient)
}

func (s *Suite) TearDownTest() {
	s.profileRepo.Items = []userEntity.UserProfile{}
	s.guildRepo.Items = []guildEntity.Guild{}
}

func (s *Suite) TearDownSubTest() {
	s.profileRepo.Items = []userEntity.UserProfile{}
	s.guildRepo.Items = []guildEntity.Guild{}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
