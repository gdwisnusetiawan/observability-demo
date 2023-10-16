package fleetTask

import (
	"context"
	"encoding/json"
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

	type message struct {
		FleetTaskId uint64 `json:"fleet-task-id"`
		Number      string `json:"number"`
	}

	msg := message{
		FleetTaskId: foundFleetTask.ID,
		Number:      foundFleetTask.Name,
	}

	msgMarshalled, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	err = BrokerProducer.ProduceMessage(ctx, "fleet-task-event", msgMarshalled)
	if err != nil {
		return nil, err
	}

	time.Sleep(4 * time.Millisecond)

	return foundFleetTask, nil
}
