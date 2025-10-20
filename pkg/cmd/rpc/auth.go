package rpc

import (
	"context"
	"os"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc/contrib/tracing"
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
	cfg := &startrpc.Config{
		ListenIP:        os.Getenv("LISTEN_IP"),
		ListenPort:      os.Getenv("LISTEN_PORT"),
		RegisterIP:      os.Getenv("REGISTER_IP"),
		RegistryIP:      os.Getenv("REGISTRY_IP"),
		RPCRegisterName: consts.AuthServiceName,
		RPCServiceVer:   consts.AuthServiceVer,
		MetricAddr:      "",
		CollectorAddr:   os.Getenv("COLLECTOR_ADDR"),
		ServerOpts:      authZrpcServerOption(),
		RPCStart:        auth.Start,
	}

	return startrpc.Start(context.Background(), cfg)
}

func authZrpcServerOption() []zrpc.ServerOption {
	return []zrpc.ServerOption{
		zrpc.WithStatsHandler(tracing.NewServerHandler()),
		zrpc.WithChainMiddleware([]zrpc.ServerMiddleware{interceptor.CtxMDInterceptor(), interceptor.ResponseInterceptor()}),
	}
}
