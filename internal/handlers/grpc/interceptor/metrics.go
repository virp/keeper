package interceptor

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func Metrics() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		h, err := handler(ctx, req)
		log.Printf("request execution time %s, method: %s", time.Since(start), info.FullMethod)
		return h, err
	}
}
