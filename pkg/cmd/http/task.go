package http

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/crazyfrankie/zrpc-todolist/interfaces/task"
	"github.com/crazyfrankie/zrpc-todolist/pkg/cmd"
	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/starthttp"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/program"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
)

type TaskCmd struct {
	*cmd.RootCmd
}

func NewTaskCmd() *TaskCmd {
	taskCmd := &TaskCmd{
		RootCmd: cmd.NewRootCmd(program.GetProcessName(), consts.TaskApiName),
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
	cfg := &starthttp.Config{
		ListenAddr:      os.Getenv("LISTEN_ADDR"),
		ServiceName:     consts.TaskApiName,
		ServiceVer:      consts.TaskApiVer,
		RegistryIP:      os.Getenv("REGISTRY_IP"),
		ShutdownTimeout: time.Second * 5,
		MetricAddr:      "",
		CollectorAddr:   os.Getenv("COLLECTOR_ADDR"),
		InitFunc:        task.Start,
	}

	return starthttp.Start(context.Background(), cfg)
}
