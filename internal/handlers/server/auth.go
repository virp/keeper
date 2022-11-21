package server

import (
	"context"

	pb "keeper/gen/service"
)

// Login implement rpc for user login call.
func (s *KeeperServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.authService.Auth(ctx, in.GetLogin(), in.GetPassword())
	if err != nil {
		return nil, err
	}

	response := pb.LoginResponse{
		Token: token,
	}
	return &response, nil
}

// Register implement rpc for user registration call.
func (s *KeeperServer) Register(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	toke, err := s.authService.Register(ctx, in.GetLogin(), in.GetPassword())
	if err != nil {
		return nil, err
	}

	response := pb.LoginResponse{
		Token: toke,
	}
	return &response, nil
}
