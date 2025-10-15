package rpc

import (
	"context"
	"os"

	"github.com/crazyfrankie/zrpc"
	"github.com/spf13/cobra"

	"github.com/crazyfrankie/zrpc-todolist/apps/user"
	"github.com/crazyfrankie/zrpc-todolist/pkg/cmd"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/program"
	"github.com/crazyfrankie/zrpc-todolist/pkg/metrics"
	"github.com/crazyfrankie/zrpc-todolist/pkg/zrpc/interceptor"
	"github.com/crazyfrankie/zrpc-todolist/pkg/zrpc/startrpc"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
)

type UserCmd struct {
	*cmd.RootCmd
}

func NewUserCmd() *UserCmd {
	userCmd := &UserCmd{
		RootCmd: cmd.NewRootCmd(program.GetProcessName(), consts.UserServiceName),
	}
	userCmd.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return userCmd.runE()
	}

	return userCmd
}

func (u *UserCmd) Exec() error {
	return u.Execute()
}

func (u *UserCmd) runE() error {
	listenIP := os.Getenv("LISTEN_IP")
	registerIP := os.Getenv("REGISTER_IP")
	listenPort := os.Getenv("LISTEN_PORT")
	metricAddr := os.Getenv("METRIC_ADDR")
	registryIP := os.Getenv("REGISTRY_IP")

	metrics.RegistryUser()

	return startrpc.Start(context.Background(), listenIP, registerIP, listenPort, metricAddr, registryIP, consts.UserServiceName, user.Start, userGrpcServerOption()...)
}

func userGrpcServerOption() []zrpc.ServerOption {
	return []zrpc.ServerOption{
		zrpc.WithChainMiddleware([]zrpc.ServerMiddleware{interceptor.CtxMDInterceptor(), interceptor.ResponseInterceptor()}),
	}
}
