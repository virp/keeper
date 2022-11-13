package server

import (
	"context"
	"errors"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "keeper/gen/service"
	"keeper/internal/entity"
	"keeper/internal/handlers/server/interceptor"
	"keeper/internal/services"
)

type AuthService interface {
	Auth(ctx context.Context, login string, password string) (string, error)
	Register(ctx context.Context, login string, password string) (string, error)
}

type TokenService interface {
	GetUser(ctx context.Context, token string) (entity.User, error)
}

type ItemService interface {
	Create(ctx context.Context, userID string, item services.Item) error
	Update(ctx context.Context, userID string, item services.Item) error
	Get(ctx context.Context, userID string, name string) (services.Item, error)
	Delete(ctx context.Context, userID string, name string) error
	List(ctx context.Context, userID string) ([]services.Item, error)
}

type KeeperServerConfig struct {
	Log   *zap.SugaredLogger
	Auth  AuthService
	Token TokenService
	Item  ItemService
}

type KeeperServer struct {
	pb.UnimplementedKeeperServiceServer
	server       *grpc.Server
	log          *zap.SugaredLogger
	authService  AuthService
	tokenService TokenService
	itemService  ItemService
}

func NewKeeperServer(cfg KeeperServerConfig) *KeeperServer {
	s := KeeperServer{
		log:          cfg.Log,
		authService:  cfg.Auth,
		tokenService: cfg.Token,
		itemService:  cfg.Item,
	}
	return &s
}

func (s *KeeperServer) Serve(listen net.Listener) error {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.Context(),
			interceptor.Logger(s.log),
			interceptor.Errors(s.log),
			interceptor.Panics(),
			interceptor.Auth(
				s.tokenService,
				map[string]bool{
					"/keeper.KeeperService/Login":    true,
					"/keeper.KeeperService/Register": true,
				},
			),
		),
	)
	pb.RegisterKeeperServiceServer(server, s)
	s.server = server
	return server.Serve(listen)
}

func (s *KeeperServer) Shutdown(ctx context.Context) error {
	wait := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(wait)
	}()

	select {
	case <-wait:
		return nil
	case <-ctx.Done():
		if _, ok := ctx.Deadline(); ok {
			return errors.New("grpc server graceful shutdown deadline exceeded")
		}
		return nil
	}
}

func (s *KeeperServer) Close() {
	s.server.Stop()
}

// TODO: Improve User ID extract logic
func getUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return ""
	}
	return userID
}
