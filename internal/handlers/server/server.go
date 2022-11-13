package server

import (
	"context"
	"net"

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

type KeeperServer struct {
	pb.UnimplementedKeeperServiceServer
	server       *grpc.Server
	authService  AuthService
	tokenService TokenService
	itemService  ItemService
}

func NewKeeperServer(authService AuthService, tokenService TokenService, itemService ItemService) *KeeperServer {
	s := KeeperServer{
		authService:  authService,
		tokenService: tokenService,
		itemService:  itemService,
	}
	return &s
}

func (s *KeeperServer) Serve(listen net.Listener) error {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.Logger(),
			interceptor.Errors(),
			interceptor.Metrics(),
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

func (s *KeeperServer) Shutdown() {
	s.server.GracefulStop()
}

// TODO: Improve User ID extract logic
func getUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return ""
	}
	return userID
}
