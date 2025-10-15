package user

import (
	"context"
	"net/http"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc/registry"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/zrpc-todolist/infra/impl/balancer"
	"github.com/crazyfrankie/zrpc-todolist/interfaces/user/handler"
	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/middleware"
	"github.com/crazyfrankie/zrpc-todolist/protocol/auth"
	"github.com/crazyfrankie/zrpc-todolist/protocol/user"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
)

// Start returns gin.Engine.
func Start(ctx context.Context, client *registry.TcpClient) (http.Handler, error) {
	srv := gin.Default()

	userService, err := client.GetService(consts.UserServiceName)
	if err != nil {
		return nil, err
	}
	authServices, err := client.GetService(consts.AuthServiceName)
	if err != nil {
		return nil, err
	}
	userCC, err := getConn(ctx, userService)
	if err != nil {
		return nil, err
	}
	authCC, err := getConn(ctx, authServices)
	if err != nil {
		return nil, err
	}

	userCli := user.NewUserServiceClient(userCC)
	authCli := auth.NewAuthServiceClient(authCC)
	userHdl := handler.NewUserHandler(userCli)
	authHdl, err := middleware.NewAuthnHandler(authCli)
	if err != nil {
		return nil, err
	}

	srv.Use(middleware.Metric(), authHdl.IgnorePath([]string{"/api/user/login", "/api/user/register"}).Auth())

	apiGroup := srv.Group("api")
	userHdl.RegisterRoute(apiGroup)

	return srv, nil
}

func getConn(ctx context.Context, services []string) (zrpc.ClientInterface, error) {
	bl := balancer.NewRoundRobinBalancer(services)
	addr, err := bl.Next(ctx)
	if err != nil {
		return nil, err
	}

	return zrpc.NewClient(addr)
}
