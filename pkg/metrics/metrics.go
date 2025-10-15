package metrics

import (
	"errors"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	path = "metrics"
)

func init() {
	registry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)
}

var registry = &prometheusRegistry{prometheus.NewRegistry()}

type prometheusRegistry struct {
	*prometheus.Registry
}

func (p *prometheusRegistry) MustRegister(cs ...prometheus.Collector) {
	for _, c := range cs {
		if err := p.Registry.Register(c); err != nil {
			if errors.As(err, &prometheus.AlreadyRegisteredError{}) {
				continue
			}
			panic(err)
		}
	}
}

func Start(listener net.Listener) error {
	mux := http.NewServeMux()
	mux.Handle(path, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	return http.Serve(listener, mux)
}
