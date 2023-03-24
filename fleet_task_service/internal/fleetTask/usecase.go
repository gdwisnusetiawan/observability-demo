package fleetTask

import (
	"context"
	"fathil/go-observability/fleet_task_service/internal/entities"
	"fathil/go-observability/pkg/observability"
	"time"
)

func IndexUseCase(ctx context.Context) (*entities.FleetTask, error) {
	ctx, span := observability.NewTraceSpan(ctx, "IndexUseCase")
	defer span.End()

	time.Sleep(2 * time.Millisecond)
	foundFleetTask, err := GetFleetTask(ctx)
	if err != nil {
		return nil, err
	}

	time.Sleep(4 * time.Millisecond)

	return foundFleetTask, nil
}
