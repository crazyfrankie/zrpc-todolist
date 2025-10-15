package signal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitExit() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	<-signals
}

func CtxWaitExit(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	select {
	case sig := <-signals:
		return fmt.Errorf("received signal: %v", sig)
	case <-ctx.Done():
		return ctx.Err()
	}
}
