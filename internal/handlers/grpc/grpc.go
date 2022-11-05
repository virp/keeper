package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	pb "keeper/gen/service"
	"keeper/internal/handlers/grpc/interceptor"
	"keeper/internal/services"
)

type ItemService interface {
	Create(ctx context.Context, userID string, item services.Item) error
	Update(ctx context.Context, userID string, item services.Item) error
	Get(ctx context.Context, userID string, name string) (services.Item, error)
	Delete(ctx context.Context, userID string, name string) error
	List(ctx context.Context, userID string) ([]services.Item, error)
}

type KeeperServer struct {
	pb.UnimplementedKeeperServiceServer
	server      *grpc.Server
	itemService ItemService
}

func NewKeeperServer(itemService ItemService) *KeeperServer {
	s := KeeperServer{
		itemService: itemService,
	}
	return &s
}

func (s *KeeperServer) Serve(listen net.Listener) error {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.Metrics(),
			interceptor.Logger(),
		),
	)
	pb.RegisterKeeperServiceServer(server, s)
	s.server = server
	return server.Serve(listen)
}

func (s *KeeperServer) Shutdown() {
	s.server.GracefulStop()
}

func getUserIDFromContext(_ context.Context) string {
	return ""
}
