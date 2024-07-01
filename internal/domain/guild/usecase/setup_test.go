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
	userRepo      *inmemory.InMemoryUserRepository
	profileRepo   *inmemory.InMemoryUserProfileRepository
	guildRepo     *inmemory.InMemoryGuildRepository
	memberRepo    *inmemory.InMemoryGuildMemberRepository
	channelRepo   *inmemory.InMemoryChannelRepository
	guildService  *GuildService
}

func (s *Suite) SetupSuite() {
	s.guildRepo = inmemory.NewInMemoryGuildRepository()
	s.profileRepo = inmemory.NewInMemoryUserProfileRepository()
	s.userRepo = inmemory.NewInMemoryUserRepository()
	s.memberRepo = inmemory.NewInMemoryGuildMemberRepository()
	s.channelRepo = inmemory.NewInMemoryChannelRepository()
	s.profileClient = inmemory.NewInMemoryUserProfileClient(s.profileRepo, s.userRepo)

	s.guildService = NewGuildService(s.guildRepo, s.memberRepo, s.channelRepo, s.profileClient)
}

func (s *Suite) TearDownTest() {
	s.profileRepo.Items = []userEntity.UserProfile{}
	s.userRepo.Items = []userEntity.User{}
	s.guildRepo.Items = []guildEntity.Guild{}
	s.memberRepo.Items = []guildEntity.GuildMember{}
	s.channelRepo.Items = []guildEntity.Channel{}
}

func (s *Suite) TearDownSubTest() {
	s.profileRepo.Items = []userEntity.UserProfile{}
	s.userRepo.Items = []userEntity.User{}
	s.memberRepo.Items = []guildEntity.GuildMember{}
	s.channelRepo.Items = []guildEntity.Channel{}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
