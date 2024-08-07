package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmingruby/telephony/internal/config"
	guildUc "github.com/charmingruby/telephony/internal/domain/guild/usecase"
	userUc "github.com/charmingruby/telephony/internal/domain/user/usecase"

	"github.com/charmingruby/telephony/internal/infra/database"
	"github.com/charmingruby/telephony/internal/infra/database/client"
	"github.com/charmingruby/telephony/internal/infra/security/cryptography"
	"github.com/charmingruby/telephony/internal/infra/security/token"
	"github.com/charmingruby/telephony/internal/infra/transport/rest"
	restEp "github.com/charmingruby/telephony/internal/infra/transport/rest/endpoint"
	"github.com/charmingruby/telephony/internal/infra/transport/ws"
	wsEp "github.com/charmingruby/telephony/internal/infra/transport/ws/endpoint"
	"github.com/charmingruby/telephony/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("CONFIGURATION: .env file not found")
	}

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error(fmt.Sprintf("CONFIGURATION: %s", err.Error()))
		os.Exit(1)
	}

	db, err := postgres.NewPostgresConnection(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE: %s", err.Error()))
		os.Exit(1)
	}

	router := gin.Default()

	initDependencies(cfg, db, router)

	server := rest.NewServer(router, cfg.ServerConfig.Port)

	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("REST SERVER: %s", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	slog.Info("HTTP Server interruption received!")

	ctx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Server.Shutdown(ctx); err != nil {
		slog.Error(fmt.Sprintf("GRACEFUL SHUTDOWN REST SERVER: %s", err.Error()))
		os.Exit(1)
	}

	slog.Info("Gracefully shutdown!")
}

func initDependencies(cfg *config.Config, db *sqlx.DB, router *gin.Engine) {
	userRepo, err := database.NewPostgresUserRepository(db)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}

	profileRepo, err := database.NewPostgresUserProfileRepository(db)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}

	guildRepo, err := database.NewPostgresGuildRepository(db)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}

	guildMemberRepo, err := database.NewPostgresGuildMemberRepository(db)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}

	channelRepo, err := database.NewPostgresChannelRepository(db)
	if err != nil {
		slog.Error(fmt.Sprintf("DATABASE REPOSITORY: %s", err.Error()))
		os.Exit(1)
	}

	userClient := client.NewUserClient(profileRepo, userRepo)
	token := token.NewJWTService(cfg.JWTConfig.SecretKey, cfg.JWTConfig.Issuer)
	crypto := cryptography.NewCryptography()

	userSvc := userUc.NewUserService(userRepo, profileRepo, crypto)
	guildSvc := guildUc.NewGuildService(guildRepo, guildMemberRepo, channelRepo, userClient)

	hub := ws.NewHub(channelRepo)
	hub.RegisterRooms()
	go hub.Start()

	restEp.NewHandler(router, token, userSvc, guildSvc).Register()
	wsEp.NewWSHandler(router, guildSvc, token, hub).Register()

}
