package rpc

import (
	"context"
	"os"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc/contrib/tracing"
	"github.com/spf13/cobra"

	"github.com/crazyfrankie/zrpc-todolist/pkg/cmd"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/program"
	"github.com/crazyfrankie/zrpc-todolist/pkg/zrpc/interceptor"
	"github.com/crazyfrankie/zrpc-todolist/pkg/zrpc/startrpc"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
)

type TaskCmd struct {
	*cmd.RootCmd
}

func NewTaskCmd() *TaskCmd {
	taskCmd := &TaskCmd{
		RootCmd: cmd.NewRootCmd(program.GetProcessName(), consts.TaskServiceName),
	}
	taskCmd.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return taskCmd.runE()
	}

	return taskCmd
}

func (u *TaskCmd) Exec() error {
	return u.Execute()
}

func (u *TaskCmd) runE() error {
	cfg := &startrpc.Config{
		ListenIP:        os.Getenv("LISTEN_IP"),
		ListenPort:      os.Getenv("LISTEN_PORT"),
		RegisterIP:      os.Getenv("REGISTER_IP"),
		RegistryIP:      os.Getenv("REGISTRY_IP"),
		RPCRegisterName: consts.TaskServiceName,
		RPCServiceVer:   consts.TaskServiceVer,
		MetricAddr:      "",
		CollectorAddr:   os.Getenv("COLLECTOR_URL"),
		ServerOpts:      taskZrpcServerOption(),
		RPCStart:        nil,
	}

	return startrpc.Start(context.Background(), cfg)
}

func taskZrpcServerOption() []zrpc.ServerOption {
	return []zrpc.ServerOption{
		zrpc.WithStatsHandler(tracing.NewServerHandler()),
		zrpc.WithChainMiddleware([]zrpc.ServerMiddleware{interceptor.CtxMDInterceptor(), interceptor.ResponseInterceptor()}),
	}
}
