package interceptor

import (
	"context"

	"github.com/crazyfrankie/zrpc"
)

func Trace() zrpc.ServerMiddleware {
	return func(ctx context.Context, req any, info *zrpc.ServerInfo, handler zrpc.Handler) (resp any, err error) {
		// TODO

		return resp, err
	}
}
