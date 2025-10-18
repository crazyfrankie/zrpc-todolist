package task

import (
	"context"
	"net/http"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc-todolist/interfaces/task/handler"
	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/middleware"
	"github.com/crazyfrankie/zrpc-todolist/protocol/auth"
	"github.com/crazyfrankie/zrpc-todolist/protocol/task"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
	"github.com/gin-gonic/gin"
)

func Start(ctx context.Context, getConn func(service string) (zrpc.ClientInterface, error), middlewares ...gin.HandlerFunc) (http.Handler, error) {
	srv := gin.Default()

	taskCC, err := getConn(consts.TaskServiceName)
	if err != nil {
		return nil, err
	}
	authCC, err := getConn(consts.AuthServiceName)
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

	middlewares = append(middlewares, authHdl.Auth())

	srv.Use(middlewares...)

	apiGroup := srv.Group("api")
	taskHdl.RegisterRoute(apiGroup)

	return srv, nil
}
