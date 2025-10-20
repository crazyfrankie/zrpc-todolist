package rpc

import (
	"context"
	"os"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc/contrib/tracing"
	"github.com/spf13/cobra"

	"github.com/crazyfrankie/zrpc-todolist/apps/user"
	"github.com/crazyfrankie/zrpc-todolist/pkg/cmd"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/program"
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
	cfg := &startrpc.Config{
		ListenIP:        os.Getenv("LISTEN_IP"),
		ListenPort:      os.Getenv("LISTEN_PORT"),
		RegisterIP:      os.Getenv("REGISTER_IP"),
		RegistryIP:      os.Getenv("REGISTRY_IP"),
		RPCRegisterName: consts.UserServiceName,
		RPCServiceVer:   consts.UserServiceVer,
		MetricAddr:      "",
		CollectorAddr:   os.Getenv("COLLECTOR_ADDR"),
		ServerOpts:      userZrpcServerOption(),
		RPCStart:        user.Start,
	}

	return startrpc.Start(context.Background(), cfg)
}

func userZrpcServerOption() []zrpc.ServerOption {
	return []zrpc.ServerOption{
		zrpc.WithStatsHandler(tracing.NewServerHandler()),
		zrpc.WithChainMiddleware([]zrpc.ServerMiddleware{interceptor.CtxMDInterceptor(), interceptor.ResponseInterceptor()}),
	}
}
