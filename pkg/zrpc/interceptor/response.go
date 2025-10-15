package interceptor

import (
	"context"
	"errors"

	"github.com/crazyfrankie/zrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/crazyfrankie/zrpc-todolist/pkg/errorx"
	"github.com/crazyfrankie/zrpc-todolist/pkg/logs"
)

func ResponseInterceptor() zrpc.ServerMiddleware {
	return func(ctx context.Context, req any, info *zrpc.ServerInfo, handler zrpc.Handler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			var customErr errorx.StatusError

			if errors.As(err, &customErr) && customErr.Code() != 0 {
				logs.CtxWarnf(ctx, "[ErrorX] error:  %v %v \n", customErr.Code(), err)
				err = status.Errorf(codes.Code(customErr.Code()), customErr.Msg())
				return
			}

			logs.CtxErrorf(ctx, "[InternalError]  error: %v \n", err)
			err = status.Errorf(codes.Internal, "internal error")
		}

		return
	}
}
