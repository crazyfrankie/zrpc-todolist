package rpc

import (
	"context"
	"os"

	"github.com/crazyfrankie/zrpc"
	"github.com/spf13/cobra"

	"github.com/crazyfrankie/zrpc-todolist/apps/auth"
	"github.com/crazyfrankie/zrpc-todolist/pkg/cmd"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/program"
	"github.com/crazyfrankie/zrpc-todolist/pkg/zrpc/interceptor"
	"github.com/crazyfrankie/zrpc-todolist/pkg/zrpc/startrpc"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
)

type AuthCmd struct {
	*cmd.RootCmd
}

func NewAuthCmd() *AuthCmd {
	authCmd := &AuthCmd{
		RootCmd: cmd.NewRootCmd(program.GetProcessName(), consts.AuthServiceName),
	}
	authCmd.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return authCmd.runE()
	}

	return authCmd
}

func (a *AuthCmd) Exec() error {
	return a.Execute()
}

func (a *AuthCmd) runE() error {
	listenIP := os.Getenv("LISTEN_IP")
	registerIP := os.Getenv("REGISTER_IP")
	listenPort := os.Getenv("LISTEN_PORT")

	return startrpc.Start(context.Background(), listenIP, registerIP, listenPort, "", consts.AuthServiceName, auth.Start, authGrpcServerOption()...)
}

func authGrpcServerOption() []zrpc.ServerOption {
	return []zrpc.ServerOption{
		zrpc.WithChainMiddleware([]zrpc.ServerMiddleware{interceptor.CtxMDInterceptor(), interceptor.ResponseInterceptor()}),
	}
}
