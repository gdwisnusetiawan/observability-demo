package order

import (
	"context"
	"encoding/json"
	"fathil/go-observability/order_service/internal/entities"
	"fathil/go-observability/pkg/observability"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
)

func GetOrder(ctx context.Context) (*entities.Order, error) {
	_, span := observability.NewTraceSpan(ctx, "GetOrder")
	defer span.End()

	time.Sleep(30 * time.Millisecond)
	return &entities.Order{
		ID:   1,
		Name: "DO-001",
	}, nil
}

func GetFleetTask(ctx context.Context) (*entities.FleetTask, error) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:3002/fleet-task", nil)
	if err != nil {
		return nil, err
	}

	otelhttptrace.Inject(ctx, req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var fleetTask entities.FleetTask
	err = json.Unmarshal(resBody, &fleetTask)
	if err != nil {
		return nil, err
	}

	return &fleetTask, nil
}
