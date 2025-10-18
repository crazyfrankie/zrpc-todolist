package startrpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc-todolist/pkg/tracing"
	zrpctracing "github.com/crazyfrankie/zrpc/contrib/tracing"
	"github.com/crazyfrankie/zrpc/registry"
	"github.com/oklog/run"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/signal"
	"github.com/crazyfrankie/zrpc-todolist/pkg/logs"
	"github.com/crazyfrankie/zrpc-todolist/pkg/metrics"
)

type Config struct {
	ListenIP        string
	ListenPort      string
	RegisterIP      string
	RegistryIP      string
	RPCRegisterName string
	RPCServiceVer   string

	MetricAddr    string
	CollectorAddr string

	ServerOpts []zrpc.ServerOption

	RPCStart func(ctx context.Context, srv zrpc.ServiceRegistrar, getConn func(service string) (zrpc.ClientInterface, error)) error
}

func Start(ctx context.Context, cfg *Config) error {
	client := registry.NewTcpClient(cfg.RegistryIP)

	g := &run.Group{}

	// Signal handler
	g.Add(func() error {
		return signal.CtxWaitExit(ctx)
	}, func(err error) {

	})

	// Prometheus metrics server
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

	// RPC server
	var (
		rpcServer *zrpc.Server
	)

	if cfg.CollectorAddr != "" {
		traceProvider, err := tracing.GetTraceProvider(cfg.RPCRegisterName, cfg.RPCServiceVer, cfg.CollectorAddr)
		if err != nil {
			return err
		}
		otel.SetTracerProvider(traceProvider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))
	}

	getConn := func(service string) (zrpc.ClientInterface, error) {
		target := fmt.Sprintf("registry:///%s", service)

		registryIP := os.Getenv("REGISTRY_IP")

		clientOptions := []zrpc.ClientOption{
			zrpc.DialWithTCPKeepAlive(15 * time.Second),
			zrpc.DialWithIdleTimeout(30 * time.Second),
			zrpc.DialWithHeartbeatInterval(30 * time.Second),
			zrpc.DialWithHeartbeatTimeout(5 * time.Second),
			zrpc.DialWithRegistryAddress(registryIP),
			zrpc.DialWithStatsHandler(zrpctracing.NewClientHandler()),
		}

		return zrpc.NewClient(target, clientOptions...)
	}

	onRegisterService := func(desc *zrpc.ServiceDesc, impl any) {
		if rpcServer != nil {
			rpcServer.RegisterService(desc, impl)
			return
		}

		rpcListenAddr := net.JoinHostPort(cfg.ListenIP, cfg.ListenPort)

		rpcServer = zrpc.NewServer(cfg.ServerOpts...)
		rpcServer.RegisterService(desc, impl)
		logs.CtxDebugf(ctx, "rpc start register, rpcRegisterName: %s, registerIP: %s, listenPort: %s", cfg.RPCRegisterName, cfg.RegisterIP, cfg.ListenPort)

		g.Add(func() error {
			// Register service
			if err := client.RegisterWithKeepAlive(cfg.RPCRegisterName, rpcListenAddr, nil, 120); err != nil {
				return fmt.Errorf("rpc register %s: %w", cfg.RPCRegisterName, err)
			}

			// Start serving
			return rpcServer.Serve("tcp", rpcListenAddr)
		}, func(err error) {
			if rpcServer != nil {
				// Graceful stop with timeout
				stopCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()

				done := make(chan struct{})
				go func() {
					rpcServer.GracefulStop()
					close(done)
				}()

				select {
				case <-done:
					logs.CtxInfof(ctx, "zRPC server stopped gracefully")
				case <-stopCtx.Done():
					logs.CtxWarnf(ctx, "zRPC server graceful stop timeout, forcing shutdown")
					rpcServer.Stop()
				}
			}
		})
	}

	if err := cfg.RPCStart(ctx, &zrpcServiceRegistrar{onRegisterService: onRegisterService}, getConn); err != nil {
		return err
	}

	// Run all services
	if err := g.Run(); err != nil {
		logs.Infof("program interrupted, %v", err)
		return err
	}

	logs.Infof("Server exited gracefully")

	return nil
}

type zrpcServiceRegistrar struct {
	onRegisterService func(desc *zrpc.ServiceDesc, impl any)
}

func (x *zrpcServiceRegistrar) RegisterService(desc *zrpc.ServiceDesc, impl any) {
	x.onRegisterService(desc, impl)
}
