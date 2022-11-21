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

// AuthService interface set requirements for auth rpc.
type AuthService interface {
	Auth(ctx context.Context, login string, password string) (string, error)
	Register(ctx context.Context, login string, password string) (string, error)
}

// TokenService interface set requirements for auth rpc.
type TokenService interface {
	GetUser(ctx context.Context, token string) (entity.User, error)
}

// ItemService interface set requirements for item rpc.
type ItemService interface {
	Create(ctx context.Context, userID string, item services.Item) error
	Update(ctx context.Context, userID string, item services.Item) error
	Get(ctx context.Context, userID string, name string) (services.Item, error)
	Delete(ctx context.Context, userID string, name string) error
	List(ctx context.Context, userID string) ([]string, error)
}

// KeeperServerConfig contains required dependencies for KeeperServer.
type KeeperServerConfig struct {
	Log   *zap.SugaredLogger
	Auth  AuthService
	Token TokenService
	Item  ItemService
}

// KeeperServer implement GRPC handlers for Keeper server.
type KeeperServer struct {
	pb.UnimplementedKeeperServiceServer
	server       *grpc.Server
	log          *zap.SugaredLogger
	authService  AuthService
	tokenService TokenService
	itemService  ItemService
}

// NewKeeperServer constructs new KeeperServer.
func NewKeeperServer(cfg KeeperServerConfig) *KeeperServer {
	s := KeeperServer{
		log:          cfg.Log,
		authService:  cfg.Auth,
		tokenService: cfg.Token,
		itemService:  cfg.Item,
	}
	return &s
}

// Serve run GRPC server.
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

// Shutdown run graceful shutdown for GRPC server.
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

// Close immediately close GRPC server.
func (s *KeeperServer) Close() {
	s.server.Stop()
}

func getUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return ""
	}
	return userID
}
