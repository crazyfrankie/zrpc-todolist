package starthttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc-todolist/pkg/gin/middleware"
	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/signal"
	"github.com/crazyfrankie/zrpc-todolist/pkg/logs"
	"github.com/crazyfrankie/zrpc-todolist/pkg/metrics"
	"github.com/crazyfrankie/zrpc-todolist/pkg/tracing"
	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	metrics.RegisterBFF()
}

func Start(ctx context.Context, listenAddr, metricAddr, collectorUrl, serviceName, serviceVersion, registryIP string, shutdownTimeout time.Duration,
	initFn func(ctx context.Context, getConn func(service string) (zrpc.ClientInterface, error),
		middlewares ...gin.HandlerFunc) (http.Handler, error)) error {
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

	getConn := func(service string) (zrpc.ClientInterface, error) {
		target := fmt.Sprintf("registry:///%s", service)

		clientOptions := []zrpc.ClientOption{
			zrpc.DialWithTCPKeepAlive(15 * time.Second),
			zrpc.DialWithIdleTimeout(30 * time.Second),
			zrpc.DialWithHeartbeatInterval(40 * time.Second),
			zrpc.DialWithHeartbeatTimeout(5 * time.Second),
			zrpc.DialWithRegistryAddress(registryIP),
		}

		return zrpc.NewClient(target, clientOptions...)
	}

	traceProvider, err := tracing.GetTraceProvider(serviceName, serviceVersion, collectorUrl)
	if err != nil {
		return err
	}

	middlewares := []gin.HandlerFunc{
		middleware.Trace(serviceName,
			otelgin.WithTracerProvider(traceProvider),
			otelgin.WithPropagators(propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			)),
		),
		middleware.Metric(),
	}

	engine, err := initFn(ctx, getConn, middlewares...)
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
