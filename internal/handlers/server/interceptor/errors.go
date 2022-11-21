package interceptor

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"keeper/internal/repository"
	"keeper/internal/services"
)

// Errors interceptor wrap errors in GRPC codes and log original error message.
func Errors(log *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			v := GetValues(ctx)
			log.Errorw("ERROR", "trace_d", v.TraceID, "message", err)

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
			case errors.Is(err, repository.ErrUserNotFound):
				wrappedError = status.Error(codes.Unauthenticated, "wrong login or password")
			case errors.Is(err, services.ErrUserInvalidPassword):
				wrappedError = status.Error(codes.Unauthenticated, "wrong login or password")
			case errors.Is(err, repository.ErrUserAlreadyExist):
				wrappedError = status.Error(codes.AlreadyExists, "user login already exist")
			case services.IsFieldErrors(err):
				wrappedError = status.Error(codes.InvalidArgument, err.Error())
			default:
				wrappedError = status.Error(codes.Internal, "internal server error")
			}
			return nil, wrappedError
		}
		return resp, nil
	}
}
