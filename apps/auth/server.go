package auth

import (
	"context"

	"github.com/crazyfrankie/zrpc"

	"github.com/crazyfrankie/zrpc-todolist/apps/auth/application"
	"github.com/crazyfrankie/zrpc-todolist/apps/auth/domain/service"
	"github.com/crazyfrankie/zrpc-todolist/protocol/auth"
)

func Start(ctx context.Context, srv zrpc.ServiceRegistrar, getConn func(service string) (zrpc.ClientInterface, error)) error {
	basic, err := application.Init(ctx)
	if err != nil {
		return err
	}
	authDomain := service.NewAuthDomain(&service.Components{
		TokenGen: basic.TokenGen,
	})
	appService := application.NewAuthApplicationService(authDomain)

	auth.RegisterAuthServiceServer(srv, appService)

	return nil
}
