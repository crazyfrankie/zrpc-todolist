package balancer

import "context"

type Balancer interface {
	Next(ctx context.Context) (string, error)
	Update(services []string)
	SetHealthCheck(checkFn func(addr string) bool)
}
