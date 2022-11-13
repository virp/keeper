package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"keeper/internal/handlers/server"
	"keeper/internal/repository/memory"
	"keeper/internal/services"
	"keeper/pkg/logger"
)

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
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		panic(err)
	}

	idGenerator := services.UuidGenerator{}
	passwordHasher := services.BCryptPasswordHasher{
		Cost: 13,
	}
	userRepository := memory.NewUserMemoryRepository()
	tokenRepository := memory.NewTokenMemoryRepository()
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
	itemRepository := memory.NewItemMemoryRepository()
	itemService := services.NewItemService(&idGenerator, itemRepository)

	s := server.NewKeeperServer(server.KeeperServerConfig{
		Log:   log,
		Auth:  authService,
		Token: tokenService,
		Item:  itemService,
	})

	serverError := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Infow("startup", "status", "grpc server started", "host", listen.Addr())
		serverError <- s.Serve(listen)
	}()

	select {
	case err := <-serverError:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			s.Close()
			return fmt.Errorf("could not stop grpc server gracefully: %w", err)
		}
	}

	return nil
}
