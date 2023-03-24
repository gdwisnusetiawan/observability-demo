package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func InitMeterProvider(appName string, endpoint string) (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(2*time.Second))),
		sdkmetric.WithResource(initResource(appName)),
	)

	global.SetMeterProvider(mp)

	return mp, nil
}
