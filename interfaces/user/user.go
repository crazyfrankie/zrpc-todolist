package user

import (
	"context"
	"net/http"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc-todolist/interfaces/user/handler"
	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/middleware"
	"github.com/crazyfrankie/zrpc-todolist/protocol/auth"
	"github.com/crazyfrankie/zrpc-todolist/protocol/user"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
	"github.com/gin-gonic/gin"
)

// Start returns gin.Engine.
func Start(ctx context.Context, getConn func(service string) (zrpc.ClientInterface, error), middlewares ...gin.HandlerFunc) (http.Handler, error) {
	srv := gin.Default()

	userCC, err := getConn(consts.UserServiceName)
	if err != nil {
		return nil, err
	}
	authCC, err := getConn(consts.AuthServiceName)
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

	middlewares = append(middlewares, authHdl.IgnorePath([]string{"/api/user/login", "/api/user/register"}).Auth())

	srv.Use(middlewares...)

	apiGroup := srv.Group("api")
	userHdl.RegisterRoute(apiGroup)

	return srv, nil
}
