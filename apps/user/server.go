package user

import (
	"context"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc-todolist/apps/user/application"
	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/repository"
	"github.com/crazyfrankie/zrpc-todolist/apps/user/domain/service"
	"github.com/crazyfrankie/zrpc-todolist/protocol/user"
)

func Start(ctx context.Context, srv zrpc.ServiceRegistrar) error {
	basic, err := application.Init(ctx)
	if err != nil {
		return err
	}
	userRepo := repository.NewUserRepository(basic.DB)
	userDomain := service.NewUserDomain(&service.Components{
		UserRepo: userRepo,
		IDGen:    basic.IDGen,
		IconOSS:  basic.IconOSS,
	})
	appService := application.NewUserApplicationService(userDomain, basic.AuthCli)

	user.RegisterUserServiceServer(srv, appService)

	return nil
}
