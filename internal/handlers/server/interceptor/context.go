package interceptor

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type ctxKey = int

const key ctxKey = 1

// Values struct store GRPC request context.
type Values struct {
	TraceID string
	Now     time.Time
}

// GetValues returns Values stored in context.
func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return &Values{
			TraceID: "00000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}
	return v
}

// Context interceptor inject additional values to GRPC request context.
func Context() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now(),
		}
		ctx = context.WithValue(ctx, key, &v)

		return handler(ctx, req)
	}
}
