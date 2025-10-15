package main

import (
	"github.com/crazyfrankie/zrpc-todolist/pkg/cmd/http"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/program"
)

func main() {
	if err := http.NewUserCmd().Exec(); err != nil {
		program.ExitWithError(err)
	}
}
