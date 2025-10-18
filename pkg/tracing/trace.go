package tracing

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
)

func GetTraceProvider(service, version, collectorUrl string) (*trace.TracerProvider, error) {
	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceVersionKey.String(version)))
	if err != nil {
		return nil, fmt.Errorf("failed create resource, %s", err)
	}

	tp, err := newTraceProvider(res, collectorUrl)
	if err != nil {
		return nil, err
	}

	return tp, nil
}

func newTraceProvider(res *resource.Resource, collectorUrl string) (*trace.TracerProvider, error) {
	exporter, err := zipkin.New(collectorUrl)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)), trace.WithResource(res))

	return traceProvider, nil
}
