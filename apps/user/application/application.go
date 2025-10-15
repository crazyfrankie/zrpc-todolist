package application

import (
	"context"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc/registry"
	"gorm.io/gorm"

	"github.com/crazyfrankie/zrpc-todolist/infra/contract/idgen"
	"github.com/crazyfrankie/zrpc-todolist/infra/contract/storage"
	"github.com/crazyfrankie/zrpc-todolist/infra/impl/balancer"
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

func Init(ctx context.Context, client *registry.Registry) (*BasicServices, error) {
	basic := &BasicServices{}
	var err error

	basic.DB, err = mysql.New()
	if err != nil {
		return nil, err
	}

	cacheCli := redis.New()

	basic.IDGen, err = idgenimpl.New(cacheCli)
	if err != nil {
		return nil, err
	}

	services, err := client.GetService(consts.AuthServiceName)
	if err != nil {
		return nil, err
	}

	authCC, err := getConn(ctx, services)
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

func getConn(ctx context.Context, services []string) (zrpc.ClientInterface, error) {
	bl := balancer.NewRoundRobinBalancer(services)
	addr, err := bl.Next(ctx)
	if err != nil {
		return nil, err
	}

	return zrpc.NewClient(addr)
}
