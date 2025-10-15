package task

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/crazyfrankie/zrpc"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/zrpc-todolist/interfaces/task/handler"
	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/middleware"
	"github.com/crazyfrankie/zrpc-todolist/protocol/auth"
	"github.com/crazyfrankie/zrpc-todolist/protocol/task"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
)

func Start(ctx context.Context) (http.Handler, error) {
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

	srv.Use(middleware.Metric(), authHdl.Auth())

	apiGroup := srv.Group("api")
	taskHdl.RegisterRoute(apiGroup)

	return srv, nil
}

func getConn(serviceName string) (zrpc.ClientInterface, error) {
	target := fmt.Sprintf("registry:///%s", serviceName)

	registryIP := os.Getenv("REGISTRY_IP")

	clientOptions := []zrpc.ClientOption{
		zrpc.DialWithTCPKeepAlive(15 * time.Second),
		zrpc.DialWithIdleTimeout(30 * time.Second),
		zrpc.DialWithHeartbeatInterval(40 * time.Second),
		zrpc.DialWithHeartbeatTimeout(5 * time.Second),
		zrpc.DialWithRegistryAddress(registryIP),
	}

	return zrpc.NewClient(target, clientOptions...)
}
