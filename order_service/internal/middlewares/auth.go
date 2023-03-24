package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

// This is the example on how to inject something to the Auth
func Auth(c *fiber.Ctx) error {
	// Any property that exists in the request | In this example just use static value
	companyName := "megaduta"

	companyMember, _ := baggage.NewMember("company", companyName)
	b, _ := baggage.New(companyMember)
	ctx := baggage.ContextWithBaggage(c.UserContext(), b)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("company", companyName))

	c.SetUserContext(ctx)

	return c.Next()
}
