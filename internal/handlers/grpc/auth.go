package grpc

import (
	"context"

	pb "keeper/gen/service"
)

func (s *KeeperServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.authService.Auth(ctx, in.Login, in.Password)
	if err != nil {
		return nil, err
	}

	response := pb.LoginResponse{
		Token: token,
	}
	return &response, nil
}
