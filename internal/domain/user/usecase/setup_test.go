package usecase

import (
	"testing"

	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/charmingruby/telephony/test/fake"
	"github.com/charmingruby/telephony/test/inmemory"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	userRepo    *inmemory.InMemoryUserRepository
	profileRepo *inmemory.InMemoryUserProfileRepository
	userService *UserService
}

func (s *Suite) SetupSuite() {
	fakeCrypto := fake.NewFakeCryptography()
	s.userRepo = inmemory.NewInMemoryUserRepository()
	s.profileRepo = inmemory.NewInMemoryUserProfileRepository()
	s.userService = NewUserService(s.userRepo, s.profileRepo, fakeCrypto)
}

func (s *Suite) SetupTest() {
	s.userRepo.Items = []entity.User{}
	s.profileRepo.Items = []entity.UserProfile{}
}

func (s *Suite) SetupSubTest() {
	s.userRepo.Items = []entity.User{}
	s.profileRepo.Items = []entity.UserProfile{}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
