package integration

import (
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/charmingruby/telephony/internal/domain/example/repository"
	exampleUc "github.com/charmingruby/telephony/internal/domain/example/usecase"
	userUc "github.com/charmingruby/telephony/internal/domain/user/usecase"
	"github.com/charmingruby/telephony/internal/infra/database"
	"github.com/charmingruby/telephony/internal/infra/transport/rest"
	"github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
	"github.com/charmingruby/telephony/test/container"
	"github.com/charmingruby/telephony/test/fake"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
	exampleRepo repository.ExampleRepository
}

func (s *Suite) SetupSuite() {
	tdb := container.NewPostgresTestDatabase()
	s.container = tdb
}

func (s *Suite) TearDownSuite() {
	s.container.DB.Close()
}

func (s *Suite) SetupTest() {
	err := s.container.RunMigrations()
	assert.NoError(s.T(), err)

	router := gin.Default()
	s.exampleRepo, err = database.NewPostgresExampleRepository(s.container.DB)
	if err != nil {
		slog.Error(fmt.Sprintf("INTEGRATION TEST, DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}

	userRepo, err := database.NewPostgresUserRepository(s.container.DB)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}

	profileRepo, err := database.NewPostgresUserProfileRepository(s.container.DB)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}

	exampleSvc := exampleUc.NewExampleService(s.exampleRepo)
	userSvc := userUc.NewUserService(userRepo, profileRepo, fake.NewFakeCryptography())

	s.handler = endpoint.NewHandler(router, exampleSvc, userSvc)
	s.handler.Register()
	server := rest.NewServer(router, "3000")

	s.server = httptest.NewServer(server.Router)
}

func (s *Suite) TearDownTest() {
	err := s.container.RollbackMigrations()
	assert.NoError(s.T(), err)

	s.server.Close()
}

func (s *Suite) Route(path string) string {
	return fmt.Sprintf("%s/api%s", s.server.URL, path)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
