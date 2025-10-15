package application

import (
	"context"

	"github.com/crazyfrankie/zrpc/registry"
	"gorm.io/gorm"

	"github.com/crazyfrankie/zrpc-todolist/infra/contract/idgen"
	"github.com/crazyfrankie/zrpc-todolist/infra/impl/cache/redis"
	idgenimpl "github.com/crazyfrankie/zrpc-todolist/infra/impl/idgen"
	"github.com/crazyfrankie/zrpc-todolist/infra/impl/mysql"
)

type BasicServices struct {
	DB    *gorm.DB
	IDGen idgen.IDGenerator
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

	return basic, nil
}
