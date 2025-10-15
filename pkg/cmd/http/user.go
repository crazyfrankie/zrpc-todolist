package http

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/crazyfrankie/zrpc-todolist/interfaces/user"
	"github.com/crazyfrankie/zrpc-todolist/pkg/cmd"
	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/starthttp"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/program"
	"github.com/crazyfrankie/zrpc-todolist/types/consts"
)

type UserCmd struct {
	*cmd.RootCmd
}

func NewUserCmd() *UserCmd {
	userCmd := &UserCmd{
		RootCmd: cmd.NewRootCmd(program.GetProcessName(), consts.UserApiName),
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
	listenAddr := os.Getenv("LISTEN_ADDR")
	metricAddr := os.Getenv("METRIC_ADDR")

	return starthttp.Start(context.Background(), listenAddr, metricAddr, user.Start, time.Second*5)
}
