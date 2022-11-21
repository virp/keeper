package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"keeper/internal/entity"
)

const (
	tokenField = "token"
)

// TokenService interface set requirements for Auth interceptor.
type TokenService interface {
	GetUser(ctx context.Context, token string) (entity.User, error)
}

// Auth interceptor implement user auth validation.
func Auth(tokenService TokenService, skips map[string]bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if skip, ok := skips[info.FullMethod]; ok && skip {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "empty metadata")
		}
		values := md.Get(tokenField)
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "token not provided")
		}
		token := values[0]
		if token == "" {
			return nil, status.Error(codes.Unauthenticated, "empty token provided")
		}

		user, err := tokenService.GetUser(ctx, token)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, "userID", user.ID)

		return handler(ctx, req)
	}
}
