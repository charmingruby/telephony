package integration

import (
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/charmingruby/telephony/internal/domain/guild/usecase"
	userUc "github.com/charmingruby/telephony/internal/domain/user/usecase"
	"github.com/charmingruby/telephony/internal/infra/database"
	"github.com/charmingruby/telephony/internal/infra/database/client"
	"github.com/charmingruby/telephony/internal/infra/security/cryptography"
	"github.com/charmingruby/telephony/internal/infra/security/token"
	"github.com/charmingruby/telephony/internal/infra/transport/rest"
	"github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
	"github.com/charmingruby/telephony/test/container"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

const (
	contentType = "application/json"
)

type Suite struct {
	suite.Suite
	container   *container.TestDatabase
	server      *httptest.Server
	handler     *endpoint.Handler
	token       *token.JWTService
	userRepo    *database.PostgresUserRepository
	profileRepo *database.PostgresUserProfileRepository
	guildRepo   *database.PostgresGuildRepository
	channelRepo *database.PostgresChannelRepository
}

func (s *Suite) SetupSuite() {
	tdb := container.NewPostgresTestDatabase()
	s.container = tdb
}

func (s *Suite) TearDownSuite() {
	s.container.DB.Close()
}

func (s *Suite) SetupSubTest() {
	s.setupDependencies()
}

func (s *Suite) TearDownSubTest() {
	err := s.container.RollbackMigrations()
	s.NoError(err)
	s.server.Close()
}

func (s *Suite) Route(path string) string {
	return fmt.Sprintf("%s/api%s", s.server.URL, path)
}

func (s *Suite) setupDependencies() {
	err := s.container.RunMigrations()
	s.NoError(err)

	router := gin.Default()

	userRepo, err := database.NewPostgresUserRepository(s.container.DB)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}
	s.userRepo = userRepo

	profileRepo, err := database.NewPostgresUserProfileRepository(s.container.DB)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}
	s.profileRepo = profileRepo

	guildRepo, err := database.NewPostgresGuildRepository(s.container.DB)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}
	s.guildRepo = guildRepo

	channelRepo, err := database.NewPostgresChannelRepository(s.container.DB)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}
	s.channelRepo = channelRepo

	userClient := client.NewUserClient(s.profileRepo, s.userRepo)

	userSvc := userUc.NewUserService(s.userRepo, s.profileRepo, cryptography.NewCryptography())
	guildSvc := usecase.NewGuildService(s.guildRepo, channelRepo, userClient)

	s.token = token.NewJWTService("secret", "telephony")

	s.handler = endpoint.NewHandler(router, s.token, userSvc, guildSvc)
	s.handler.Register()
	server := rest.NewServer(router, "3000")

	s.server = httptest.NewServer(server.Router)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
