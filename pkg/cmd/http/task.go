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
		RootCmd: cmd.NewRootCmd(program.GetProcessName(), consts.UserApiName),
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
	listenAddr := os.Getenv("LISTEN_ADDR")
	metricAddr := os.Getenv("METRIC_ADDR")

	return starthttp.Start(context.Background(), listenAddr, metricAddr, task.Start, time.Second*5)
}
