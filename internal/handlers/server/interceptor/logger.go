package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Logger interceptor log all GRPC requests.
func Logger(log *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		v := GetValues(ctx)
		var clientIP string
		if p, ok := peer.FromContext(ctx); ok {
			clientIP = p.Addr.String()
		}
		log.Infow(
			"request started",
			"trace_id", v.TraceID,
			"method", info.FullMethod,
			"client", clientIP,
		)
		h, err := handler(ctx, req)
		code := "OK"
		st, ok := status.FromError(err)
		if ok {
			code = st.Code().String()
		}
		log.Infow(
			"request completed",
			"trace_id", v.TraceID,
			"method", info.FullMethod,
			"client", clientIP,
			"code", code,
			"since", time.Since(v.Now),
		)
		return h, err
	}
}
