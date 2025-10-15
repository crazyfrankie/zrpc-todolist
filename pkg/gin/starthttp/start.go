package starthttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/oklog/run"

	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/signal"
	"github.com/crazyfrankie/zrpc-todolist/pkg/logs"
	"github.com/crazyfrankie/zrpc-todolist/pkg/metrics"
)

func init() {
	metrics.RegisterBFF()
}

func Start(ctx context.Context, listenAddr, metricAddr string, initFn func(ctx context.Context) (http.Handler, error), shutdownTimeout time.Duration) error {

	g := &run.Group{}

	// Signal handler
	g.Add(func() error {
		return signal.CtxWaitExit(ctx)
	}, func(err error) {

	})

	if metricAddr != "" {
		g.Add(func() error {
			listener, err := net.Listen("tcp", metricAddr)
			if err != nil {
				return err
			}

			return metrics.Start(listener)
		}, func(err error) {

		})
	}

	engine, err := initFn(ctx)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: engine,
	}

	g.Add(func() error {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server failed: %w", err)
		}
		return nil
	}, func(err error) {
		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logs.Errorf("failed to shutdown main server: %v", err)
		}
		logs.Infof("Server shutdown successfully")
	})

	if err := g.Run(); err != nil {
		logs.Infof("program interrupted, %v", err)
		return err
	}

	logs.Infof("Server exited gracefully")

	return nil
}
