package app

import (
	"context"
	"fathil/go-observability/fleet_task_service/config"
	"fathil/go-observability/fleet_task_service/internal/fleetTask"
	"fathil/go-observability/pkg/observability"
	"fmt"
	"log"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
)

type app struct {
	cfg   *config.Config
	fiber *fiber.App
}

func New(cfg *config.Config) *app {
	return &app{
		cfg:   cfg,
		fiber: fiber.New(),
	}
}

func (s *app) Run() error {
	tp := observability.InitTracerProvider(s.cfg.App.Name, s.cfg.Observability.OtelEndpoint)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	fleetTask.InitBrokerProducer()

	s.fiber.Use(otelfiber.Middleware())

	v1 := s.fiber.Group("fleet-task")
	v1.Get("", fleetTask.IndexHandler)

	return s.fiber.Listen(fmt.Sprintf(":%s", s.cfg.App.Port))
}

// func initResource(cfg *config.Config) *sdkresource.Resource {
// 	extraResource, _ := sdkresource.New(
// 		context.Background(),
// 		sdkresource.WithOS(),
// 		sdkresource.WithProcess(),
// 		sdkresource.WithContainer(),
// 		sdkresource.WithHost(),
// 		sdkresource.WithAttributes(
// 			semconv.ServiceName(cfg.App.Name),
// 		),
// 	)

// 	resource, _ := sdkresource.Merge(
// 		sdkresource.Default(),
// 		extraResource,
// 	)

// 	return resource
// }

// func initTracerProvider(cfg *config.Config) *sdktrace.TracerProvider {
// 	ctx := context.Background()

// 	exporter, err := otlptracegrpc.New(
// 		ctx,
// 		otlptracegrpc.WithInsecure(),
// 		otlptracegrpc.WithEndpoint(cfg.Observability.OtelEndpoint),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	tp := sdktrace.NewTracerProvider(
// 		sdktrace.WithBatcher(exporter),
// 		sdktrace.WithResource(initResource(cfg)),
// 	)
// 	otel.SetTracerProvider(tp)
// 	return tp
// }
