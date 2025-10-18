package starthttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/crazyfrankie/zrpc"
	zrpctracing "github.com/crazyfrankie/zrpc/contrib/tracing"
	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/middleware"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/signal"
	"github.com/crazyfrankie/zrpc-todolist/pkg/logs"
	"github.com/crazyfrankie/zrpc-todolist/pkg/metrics"
	"github.com/crazyfrankie/zrpc-todolist/pkg/tracing"
)

func init() {
	metrics.RegisterBFF()
}

type Config struct {
	ListenAddr      string
	ServiceName     string
	ServiceVer      string
	RegistryIP      string
	ShutdownTimeout time.Duration

	MetricAddr   string
	CollectorURL string

	InitFunc func(ctx context.Context, getConn func(service string) (zrpc.ClientInterface, error), middlewares ...gin.HandlerFunc) (http.Handler, error)
}

func Start(ctx context.Context, cfg *Config) error {
	g := &run.Group{}

	// Signal handler
	g.Add(func() error {
		return signal.CtxWaitExit(ctx)
	}, func(err error) {

	})

	if cfg.MetricAddr != "" {
		g.Add(func() error {
			listener, err := net.Listen("tcp", cfg.MetricAddr)
			if err != nil {
				return err
			}

			return metrics.Start(listener)
		}, func(err error) {

		})
	}

	getConn := func(service string) (zrpc.ClientInterface, error) {
		target := fmt.Sprintf("registry:///%s", service)

		clientOptions := []zrpc.ClientOption{
			zrpc.DialWithTCPKeepAlive(15 * time.Second),
			zrpc.DialWithIdleTimeout(30 * time.Second),
			zrpc.DialWithHeartbeatInterval(40 * time.Second),
			zrpc.DialWithHeartbeatTimeout(5 * time.Second),
			zrpc.DialWithRegistryAddress(cfg.RegistryIP),
			zrpc.DialWithStatsHandler(zrpctracing.NewClientHandler()),
		}

		return zrpc.NewClient(target, clientOptions...)
	}

	middlewares := []gin.HandlerFunc{
		middleware.Metric(),
	}

	if cfg.CollectorURL != "" {
		traceProvider, err := tracing.GetTraceProvider(cfg.ServiceName, cfg.ServiceVer, cfg.CollectorURL)
		if err != nil {
			return err
		}

		otel.SetTracerProvider(traceProvider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		middlewares = append([]gin.HandlerFunc{middleware.Trace(cfg.ServiceName)}, middlewares...)
	}

	engine, err := cfg.InitFunc(ctx, getConn, middlewares...)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: engine,
	}

	g.Add(func() error {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server failed: %w", err)
		}
		return nil
	}, func(err error) {
		shutdownCtx, cancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
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
