package main

import (
	"context"
	"encoding/json"
	"errors"
	"fathil/go-observability/order_service/config"
	"fathil/go-observability/pkg/observability"
	"fathil/go-observability/pkg/otelsarama"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	tp := observability.InitTracerProvider(cfg.App.Name, cfg.Observability.OtelEndpoint)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	mp, err := observability.InitMeterProvider(cfg.App.Name, cfg.Observability.OtelEndpoint)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down meter provider: %v", err)
		}
	}()

	err = NewConsumer(context.Background())
	if err != nil {
		panic(err)
	}

}

type BrokerConsumer struct {
	ready chan bool
}

func NewConsumer(ctx context.Context) error {
	log.Println("Starting a new Sarama consumer")

	version, err := sarama.ParseKafkaVersion("0.11.0.0")
	if err != nil {
		panic(err)
	}

	saramCfg := sarama.NewConfig()
	saramCfg.Version = version
	saramCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramCfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}

	cg, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "1", saramCfg)
	if err != nil {
		return err
	}

	handler := NewHandler()
	cgh := otelsarama.WrapConsumerGroupHandler(
		handler,
		otelsarama.WithTracerProvider(otel.GetTracerProvider()),
		otelsarama.WithPropagators(otel.GetTextMapPropagator()),
	)
	shutdownErr := make(chan error)
	defer close(shutdownErr)

	go gracefulShutDown(shutdownErr, cg)

	go func() {
		for {
			if err := cg.Consume(ctx, []string{"fleet-task-event"}, cgh); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error consumer: %v", err)
			}
		}
	}()

	return <-shutdownErr
}

func gracefulShutDown(shutdownError chan<- error, cg sarama.ConsumerGroup) {
	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownError <- cg.Close()
}

func NewHandler() *BrokerConsumer {
	return &BrokerConsumer{
		ready: make(chan bool),
	}
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (*BrokerConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (*BrokerConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			err := handleMsg(msg)
			if err != nil {
				return err
			}

			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

// Setup implements sarama.ConsumerGroupHandler.
func (b *BrokerConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(b.ready)
	return nil
}

func handleMsg(msg *sarama.ConsumerMessage) error {
	var message struct {
		FleetTaskId uint64 `json:"fleet-task-id"`
		Number      string `json:"number"`
	}
	err := json.Unmarshal((msg.Value), &message)
	if err != nil {
		return err
	}
	log.Println(string(msg.Value))
	return nil
}
