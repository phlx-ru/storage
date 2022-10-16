package auth

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	baseGRPC "google.golang.org/grpc"
)

func Default(ctx context.Context, endpoint string, timeout time.Duration) (*baseGRPC.ClientConn, error) {
	opts := []grpc.ClientOption{
		grpc.WithEndpoint(endpoint),
		grpc.WithTimeout(timeout),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	}
	return grpc.DialInsecure(ctx, opts...)
}
