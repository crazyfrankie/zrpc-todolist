package interceptor

import (
	"context"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc/metadata"

	"github.com/crazyfrankie/zrpc-todolist/pkg/ctxcache"
)

func CtxMDInterceptor() zrpc.ServerMiddleware {
	return func(ctx context.Context, req any, info *zrpc.ServerInfo, handler zrpc.Handler) (resp any, err error) {
		ctx = ctxcache.Init(ctx)

		md, _ := metadata.FromInComingContext(ctx)

		for k, v := range md {
			ctxcache.Store(ctx, k, v)
		}

		return handler(ctx, req)
	}
}
