package otel

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func SetupOtel(ctx context.Context, host, name string) (func() error, error) {
	res, err := resource.New(ctx,
		// resource.WithFromEnv(),
		resource.WithAttributes(
			semconv.ServiceName(name),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating resource")
	}

	traceExporterShutdown, err := setupOtelTracerProvider(ctx, res, host)
	if err != nil {
		return nil, errors.Wrap(err, "error setting up OTLP traces exporter")
	}
	metricExporterShutdown, err := setupOtelMeterProvider(ctx, res, host)
	if err != nil {
		traceExporterShutdown(ctx)
		return nil, errors.Wrap(err, "error setting up OTLP metrics exporter")
	}

	return func() error {
		err := traceExporterShutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "error shutting down OTLP traces exporter")
		}
		err = metricExporterShutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "error shutting down OTLP metrics exporter")
		}
		log.Println("Shutdown otel")
		return nil
	}, nil
}

func GetTracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("")
}

func GetMeter() metric.Meter {
	return otel.GetMeterProvider().Meter("")
}
