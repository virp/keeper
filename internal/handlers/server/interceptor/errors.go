package interceptor

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"keeper/internal/repository"
)

func Errors() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			log.Printf("request error: %s, method: %s", err, info.FullMethod)

			_, isStatusError := status.FromError(err)
			var wrappedError error
			switch {
			case isStatusError:
				wrappedError = err
			case errors.Is(err, repository.ErrItemAlreadyExist):
				wrappedError = status.Error(codes.AlreadyExists, err.Error())
			case errors.Is(err, repository.ErrItemNotFound):
				wrappedError = status.Error(codes.NotFound, err.Error())
			case errors.Is(err, repository.ErrTokenNotFound):
				wrappedError = status.Error(codes.Unauthenticated, "authentication required")
			case errors.Is(err, repository.ErrTokenExpired):
				wrappedError = status.Error(codes.Unauthenticated, "authentication required")
			default:
				wrappedError = status.Error(codes.Internal, "internal server error")
			}
			return nil, wrappedError
		}
		return resp, nil
	}
}
