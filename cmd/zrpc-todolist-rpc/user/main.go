package main

import (
	"github.com/crazyfrankie/zrpc-todolist/pkg/cmd/rpc"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/program"
)

func main() {
	if err := rpc.NewTaskCmd().Exec(); err != nil {
		program.ExitWithError(err)
	}
}
