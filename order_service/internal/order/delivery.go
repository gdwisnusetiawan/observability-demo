package order

import (
	"fathil/go-observability/pkg/observability"
	"time"

	"github.com/gofiber/fiber/v2"
)

func IndexHandler(c *fiber.Ctx) error {
	ctx, span := observability.NewTraceSpan(
		c.UserContext(),
		"IndexHandler",
	)
	defer span.End()

	time.Sleep(5 * time.Millisecond)
	foundOrder, err := IndexUseCase(ctx)
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Millisecond)

	return c.JSON(foundOrder)
}
