package startrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/crazyfrankie/zrpc"
	"github.com/crazyfrankie/zrpc/registry"
	"github.com/oklog/run"

	"github.com/crazyfrankie/zrpc-todolist/pkg/lang/signal"
	"github.com/crazyfrankie/zrpc-todolist/pkg/logs"
	"github.com/crazyfrankie/zrpc-todolist/pkg/metrics"
)

func Start(ctx context.Context, listenIP, registerIP, listenPort, metricAddr, registryIP, rpcRegisterName string,
	rpcStart func(ctx context.Context, srv zrpc.ServiceRegistrar) error,
	opts ...zrpc.ServerOption) error {

	client := registry.NewTcpClient(registryIP)

	g := &run.Group{}

	// Signal handler
	g.Add(func() error {
		return signal.CtxWaitExit(ctx)
	}, func(err error) {

	})

	// Prometheus metrics server
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

	// RPC server
	var (
		rpcServer *zrpc.Server
	)

	onRegisterService := func(desc *zrpc.ServiceDesc, impl any) {
		if rpcServer != nil {
			rpcServer.RegisterService(desc, impl)
			return
		}

		rpcListenAddr := net.JoinHostPort(listenIP, listenPort)

		rpcServer = zrpc.NewServer(opts...)
		rpcServer.RegisterService(desc, impl)
		logs.CtxDebugf(ctx, "rpc start register, rpcRegisterName: %s, registerIP: %s, listenPort: %s", rpcRegisterName, registerIP, listenPort)

		g.Add(func() error {
			// Register service
			if err := client.RegisterWithKeepAlive(rpcRegisterName, rpcListenAddr, nil, 120); err != nil {
				return fmt.Errorf("rpc register %s: %w", rpcRegisterName, err)
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

	if err := rpcStart(ctx, &zrpcServiceRegistrar{onRegisterService: onRegisterService}); err != nil {
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
