package interceptor

import (
	"context"
	"fmt"
	"runtime/debug"

	"google.golang.org/grpc"
)

func Panics() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				trace := debug.Stack()
				err = fmt.Errorf("PANIC [%v] TRACE [%s]", rec, string(trace))
			}
		}()
		return handler(ctx, req)
	}
}
