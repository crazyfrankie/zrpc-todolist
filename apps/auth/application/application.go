package application

import (
	"context"

	"github.com/crazyfrankie/zrpc-todolist/infra/contract/token"
	"github.com/crazyfrankie/zrpc-todolist/infra/impl/cache/redis"
	tokenimpl "github.com/crazyfrankie/zrpc-todolist/infra/impl/token"
)

type BasicServices struct {
	TokenGen token.Token
}

func Init(ctx context.Context) (*BasicServices, error) {
	basic := &BasicServices{}
	var err error

	cacheCli := redis.New()

	basic.TokenGen, err = tokenimpl.New(cacheCli)
	if err != nil {
		return nil, err
	}

	return basic, nil
}
