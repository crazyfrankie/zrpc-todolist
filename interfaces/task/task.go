package task

import (
	"context"
	"net/http"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc/registry"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/zrpc-todolist/infra/impl/balancer"
	"github.com/crazyfrankie/zrpc-todolist/interfaces/task/handler"
	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/middleware"
	"github.com/crazyfrankie/zrpc-todolist/protocol/auth"
	"github.com/crazyfrankie/zrpc-todolist/protocol/task"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
)

func Start(ctx context.Context, client *registry.Registry) (http.Handler, error) {
	srv := gin.Default()

	taskServices, err := client.GetService(consts.AuthServiceName)
	if err != nil {
		return nil, err
	}
	authServices, err := client.GetService(consts.AuthServiceName)
	if err != nil {
		return nil, err
	}
	taskCC, err := getConn(ctx, taskServices)
	if err != nil {
		return nil, err
	}
	authCC, err := getConn(ctx, authServices)
	if err != nil {
		return nil, err
	}

	taskCli := task.NewTaskServiceClient(taskCC)
	authCli := auth.NewAuthServiceClient(authCC)
	taskHdl := handler.NewTaskHandler(taskCli)
	authHdl, err := middleware.NewAuthnHandler(authCli)
	if err != nil {
		return nil, err
	}

	srv.Use(middleware.Metric(), authHdl.Auth())

	apiGroup := srv.Group("api")
	taskHdl.RegisterRoute(apiGroup)

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
