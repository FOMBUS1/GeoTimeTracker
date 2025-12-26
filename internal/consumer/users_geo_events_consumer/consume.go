package usersgeoeventsconsumer

import (
	"context"
	"log/slog"
	"time"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func (c *UsersGeoEventsConsumer) Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBroker,
		GroupID:           "GeoService_group",
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
		MaxWait:           500 * time.Millisecond,
	})
	defer r.Close()

	slog.Info("UsersGeoEventsConsumer started", "topic", c.topicName)

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			slog.Error("UsersGeoEventsConsumer.consume error", "error", err.Error())
		}
		slog.Info("Message got!")
		var geoInfo models.GeoKafkaMessage
		//slog.Info("Message: ", msg.Value)
		err = proto.Unmarshal(msg.Value, &geoInfo)
		if err != nil {
			slog.Error("parce", "error", err)
			continue
		}
		err = c.usersGeoEventsProcessor.Handle(ctx, &geoInfo)
		if err != nil {
			slog.Error("Handle", "error", err)
		}
		slog.Info("Message handled!")
	}

}
