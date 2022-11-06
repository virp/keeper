package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func Logger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var clientIP string
		if p, ok := peer.FromContext(ctx); ok {
			clientIP = p.Addr.String()
		}
		log.Printf("request started, method: %s, client: %s", info.FullMethod, clientIP)
		h, err := handler(ctx, req)
		code := "OK"
		st, ok := status.FromError(err)
		if ok {
			code = st.Code().String()
		}
		log.Printf("request completed, code: %s, method: %s, client: %s", code, info.FullMethod, clientIP)
		return h, err
	}
}
