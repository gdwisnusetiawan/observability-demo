package fleetTask

import (
	"context"
	"fathil/go-observability/fleet_task_service/internal/entities"
	"fathil/go-observability/pkg/observability"
	"fathil/go-observability/pkg/otelsarama"
	"time"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
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

type brokerProducer struct {
	conn sarama.SyncProducer
}

var BrokerProducer *brokerProducer

func InitBrokerProducer() {
	version, err := sarama.ParseKafkaVersion("0.11.0.0")
	if err != nil {
		panic(err)
	}

	saramCfg := sarama.NewConfig()
	saramCfg.Version = version
	saramCfg.Producer.Retry.Max = 10
	saramCfg.Producer.Return.Successes = true
	saramCfg.Producer.Return.Errors = true

	ap, err := sarama.NewSyncProducer([]string{"localhost:9092"}, saramCfg)
	if err != nil {
		panic(err)
	}

	ap = otelsarama.WrapSyncProducer(
		saramCfg,
		ap,
		otelsarama.WithTracerProvider(otel.GetTracerProvider()),
		otelsarama.WithPropagators(otel.GetTextMapPropagator()),
	)
	if err != nil {
		panic(err)
	}
	BrokerProducer = &brokerProducer{
		conn: ap,
	}
}

func (b *brokerProducer) ProduceMessage(ctx context.Context, topic string, msg []byte) error {
	message := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	}

	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(&message))

	_, _, err := b.conn.SendMessage(&message)

	return err

}
