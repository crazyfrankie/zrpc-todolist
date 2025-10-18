package application

import (
	"context"

	"github.com/crazyfrankie/zrpc"
	"gorm.io/gorm"

	"github.com/crazyfrankie/zrpc-todolist/infra/contract/idgen"
	"github.com/crazyfrankie/zrpc-todolist/infra/contract/storage"
	"github.com/crazyfrankie/zrpc-todolist/infra/impl/cache/redis"
	idgenimpl "github.com/crazyfrankie/zrpc-todolist/infra/impl/idgen"
	"github.com/crazyfrankie/zrpc-todolist/infra/impl/mysql"
	storageimpl "github.com/crazyfrankie/zrpc-todolist/infra/impl/storage"
	"github.com/crazyfrankie/zrpc-todolist/protocol/auth"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
)

type BasicServices struct {
	DB      *gorm.DB
	IDGen   idgen.IDGenerator
	IconOSS storage.Storage
	AuthCli auth.AuthServiceClient
}

func Init(ctx context.Context, getConn func(service string) (zrpc.ClientInterface, error)) (*BasicServices, error) {
	basic := &BasicServices{}
	var err error

	basic.DB, err = mysql.New()
	if err != nil {
		return nil, err
	}

	cacheCli := redis.New()

	basic.IDGen, err = idgenimpl.New(cacheCli, consts.UserServiceName)
	if err != nil {
		return nil, err
	}

	authCC, err := getConn(consts.AuthServiceName)
	if err != nil {
		return nil, err
	}

	basic.AuthCli = auth.NewAuthServiceClient(authCC)

	basic.IconOSS, err = storageimpl.New(ctx)
	if err != nil {
		return nil, err
	}

	return basic, nil
}
