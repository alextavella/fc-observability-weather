package otel

import (
	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func setupOtelMeterProvider(
	ctx context.Context,
	res *resource.Resource,
	host string,
) (func(ctx context.Context) error, error) {
	metricOpts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithInsecure(), otlpmetricgrpc.WithEndpoint(host)}
	metricExporter, err := otlpmetricgrpc.New(ctx, metricOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "createing new metric exporter failed")
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(1))),
	)
	otel.SetMeterProvider(mp)
	return mp.Shutdown, nil
}
