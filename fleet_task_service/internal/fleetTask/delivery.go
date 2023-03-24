package fleetTask

import (
	"fathil/go-observability/pkg/observability"
	"time"

	"github.com/gofiber/fiber/v2"
)

func IndexHandler(c *fiber.Ctx) error {
	ctx, span := observability.NewTraceSpan(c.UserContext(), "IndexHandler")
	defer span.End()

	time.Sleep(6 * time.Millisecond)
	foundFleetTask, err := IndexUseCase(ctx)
	if err != nil {
		return err
	}
	time.Sleep(7 * time.Millisecond)

	return c.JSON(foundFleetTask)
}
