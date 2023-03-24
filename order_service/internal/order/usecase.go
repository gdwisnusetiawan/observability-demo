package order

import (
	"context"
	"fathil/go-observability/order_service/internal/entities"
	"fathil/go-observability/pkg/observability"
	"time"
)

func IndexUseCase(ctx context.Context) (*entities.Order, error) {
	ctx, span := observability.NewTraceSpan(ctx, "IndexUseCase")
	defer span.End()

	time.Sleep(10 * time.Millisecond)

	foundOrder, err := GetOrder(ctx)
	if err != nil {
		return nil, err
	}

	foundFleetTask, err := GetFleetTask(ctx)
	if err != nil {
		return nil, err
	}

	foundOrder.FleetTask = foundFleetTask

	return foundOrder, nil
}
