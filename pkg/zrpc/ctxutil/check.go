package ctxutil

import (
	"context"

	"github.com/crazyfrankie/zrpc-todolist/pkg/errorx"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/conv"
	"github.com/crazyfrankie/zrpc-todolist/types/errno"
)

func CheckAccess(ctx context.Context, ownerUserID int64) error {
	if MustGetUserIDFromCtx(ctx) == ownerUserID {
		return nil
	}

	return errorx.New(errno.ErrNoPermissionCode, errorx.KV("ownerID", conv.Int64ToStr(ownerUserID)))
}
