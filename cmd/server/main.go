package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"keeper/internal/handlers/grpc"
	"keeper/internal/repository/memory"
	"keeper/internal/services"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
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
	s := grpc.NewKeeperServer(authService, tokenService, itemService)

	serverError := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("grpc servert started at: %s", listen.Addr())
		serverError <- s.Serve(listen)
	}()

	select {
	case err := <-serverError:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Printf("shutdown started, signal: %s", sig)
		s.Shutdown()
		log.Printf("shutdown completed, signal: %s", sig)
	}

	return nil
}
