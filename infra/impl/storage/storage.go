package storage

import (
	"context"
	"os"

	"github.com/crazyfrankie/zrpc-todolist/infra/contract/storage"
	"github.com/crazyfrankie/zrpc-todolist/infra/impl/storage/minio"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
)

type Storage = storage.Storage

func New(ctx context.Context) (Storage, error) {
	return minio.New(
		ctx,
		os.Getenv(consts.MinIOEndpoint),
		os.Getenv(consts.MinIOAK),
		os.Getenv(consts.MinIOSK),
		os.Getenv(consts.StorageBucket),
		false,
	)
}
