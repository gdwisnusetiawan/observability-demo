package fleetTask

import (
	"context"
	"fathil/go-observability/fleet_task_service/internal/entities"
	"fathil/go-observability/pkg/observability"
	"time"
)

func GetFleetTask(ctx context.Context) (*entities.FleetTask, error) {
	_, span := observability.NewTraceSpan(ctx, "GetFleetTask")
	defer span.End()

	time.Sleep(4 * time.Millisecond)
	return &entities.FleetTask{
		ID:   1,
		Name: "FO-001",
		Vehicle: entities.Vehicle{
			ID:           2,
			LisencePlate: "L 123 HK",
		},
	}, nil
}
