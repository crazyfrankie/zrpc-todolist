package balancer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNoAvailableService = errors.New("no available service")
)

type RoundRobinBalancer struct {
	services    []string
	index       int
	mu          sync.Mutex
	healthCheck func(addr string) bool
}

func NewRoundRobinBalancer(services []string) *RoundRobinBalancer {
	return &RoundRobinBalancer{
		services: services,
		index:    0,
	}
}

func (rr *RoundRobinBalancer) Next(ctx context.Context) (string, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if len(rr.services) == 0 {
		return "", ErrNoAvailableService
	}

	for i := 0; i < len(rr.services); i++ {
		service := rr.services[rr.index]
		rr.index = (rr.index + 1) % len(rr.services)

		if rr.healthCheck == nil {
			return service, nil
		}

		if rr.healthCheck(service) {
			return service, nil
		}
	}

	return "", ErrNoAvailableService
}

func (rr *RoundRobinBalancer) Update(services []string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	rr.services = services
	rr.index = 0
}

func (rr *RoundRobinBalancer) SetHealthCheck(checkFn func(addr string) bool) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	rr.healthCheck = checkFn
}
