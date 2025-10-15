package application

import (
	"context"
	"fmt"
	"os"
	"time"

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

func Init(ctx context.Context) (*BasicServices, error) {
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

func getConn(serviceName string) (zrpc.ClientInterface, error) {
	target := fmt.Sprintf("registry:///%s", serviceName)

	registryIP := os.Getenv("REGISTRY_IP")

	clientOptions := []zrpc.ClientOption{
		zrpc.DialWithTCPKeepAlive(15 * time.Second),
		zrpc.DialWithIdleTimeout(30 * time.Second),
		zrpc.DialWithHeartbeatInterval(30 * time.Second),
		zrpc.DialWithHeartbeatTimeout(5 * time.Second),
		zrpc.DialWithRegistryAddress(registryIP),
	}

	return zrpc.NewClient(target, clientOptions...)
}
