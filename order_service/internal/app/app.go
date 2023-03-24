package app

import (
	"context"
	"fathil/go-observability/order_service/config"
	"fathil/go-observability/order_service/internal/middlewares"
	"fathil/go-observability/order_service/internal/order"
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

func (a *app) Run() error {
	tp := observability.InitTracerProvider(a.cfg.App.Name, a.cfg.Observability.OtelEndpoint)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	mp, _ := observability.InitMeterProvider(a.cfg.App.Name, a.cfg.Observability.OtelEndpoint)
	defer func() {
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down meter provider: %v", err)
		}
	}()

	a.fiber.Use(otelfiber.Middleware())
	a.fiber.Use(middlewares.Auth)

	v1 := a.fiber.Group("order")
	v1.Get("", order.IndexHandler)

	return a.fiber.Listen(fmt.Sprintf(":%s", a.cfg.App.Port))
}
