package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"go.uber.org/zap"

	"keeper/internal/handlers/server"
	"keeper/internal/repository/cloud"
	"keeper/internal/repository/memory"
	"keeper/internal/services"
	s3 "keeper/pkg/cloud"
	"keeper/pkg/logger"
)

var build = "develop"

const (
	passwordHashCost = 13
)

type config struct {
	conf.Version
	Address         string        `conf:"default:0.0.0.0:3200"`
	ShutdownTimeout time.Duration `conf:"default:5s"`
	TokenLifetime   time.Duration `conf:"default:24h"`
	StorageType     string        `conf:"default:memory,help:Storage type can be memory or cloud"`
	CloudStorage    struct {
		Bucket   string `conf:"default:keeper"`
		Endpoint string `conf:"default:http://localhost:4566"`
		Region   string `conf:"default:us-east-1"`
		Key      string `conf:"default:key"`
		Secret   string `conf:"default:secret"`
	}
}

func main() {
	log, err := logger.New("KEEPER-SERVER")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		//goland:noinspection GoUnhandledErrorResult
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// =========================================================================
	// Configuration

	cfg := config{
		Version: conf.Version{
			Build: build,
			Desc:  "GophKeeper GRPC Server",
		},
	}
	const prefix = "KEEPER"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parse config: %w", err)
	}

	// =========================================================================
	// Create GRPC Server

	serverError := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	s, err := buildServer(log, cfg)
	if err != nil {
		return fmt.Errorf("build server: %w", err)
	}

	// =========================================================================
	// App Starting

	log.Infow("starting service", "version", build)

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	listen, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Infow("startup", "status", "grpc server started", "host", listen.Addr())
		serverError <- s.Serve(listen)
	}()

	// =========================================================================
	// Shutdown

	select {
	case err := <-serverError:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			s.Close()
			return fmt.Errorf("could not stop grpc server gracefully: %w", err)
		}
	}

	return nil
}

func buildServer(log *zap.SugaredLogger, cfg config) (*server.KeeperServer, error) {
	idGenerator := services.UuidGenerator{}
	passwordHasher := services.BCryptPasswordHasher{
		Cost: passwordHashCost,
	}

	var userRepository services.UserRepository
	var tokenRepository services.TokenRepository
	var itemRepository services.ItemRepository

	switch cfg.StorageType {
	case "memory":
		userRepository = memory.NewUserRepository()
		tokenRepository = memory.NewTokenRepository(cfg.TokenLifetime)
		itemRepository = memory.NewItemRepository()
	case "cloud":
		s3Client := s3.NewS3Client(
			cfg.CloudStorage.Key,
			cfg.CloudStorage.Secret,
			cfg.CloudStorage.Region,
			cfg.CloudStorage.Endpoint,
		)
		userRepository = cloud.NewUserRepository(s3Client, cfg.CloudStorage.Bucket)
		tokenRepository = cloud.NewTokenRepository(s3Client, cfg.CloudStorage.Bucket, cfg.TokenLifetime)
		itemRepository = cloud.NewItemRepository(s3Client, cfg.CloudStorage.Bucket)
	default:
		return nil, fmt.Errorf("not supported storage type: %s", cfg.StorageType)
	}

	authService := services.NewAuthService(
		&idGenerator,
		&passwordHasher,
		userRepository,
		tokenRepository,
	)
	tokenService := services.NewTokenService(
		userRepository,
		tokenRepository,
	)
	itemService := services.NewItemService(&idGenerator, itemRepository)

	s := server.NewKeeperServer(server.KeeperServerConfig{
		Log:   log,
		Auth:  authService,
		Token: tokenService,
		Item:  itemService,
	})

	return s, nil
}
