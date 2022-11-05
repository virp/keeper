package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func Logger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("request started, method: %s", info.FullMethod)
		h, err := handler(ctx, req)
		log.Printf("request completed, method: %s", info.FullMethod)
		return h, err
	}
}
