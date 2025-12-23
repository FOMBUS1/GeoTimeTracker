package kafka

import (
	"context"
	"fmt"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(w *kafka.Writer) *Producer {
	return &Producer{writer: w}
}

func (p *Producer) Send(ctx context.Context, msg *models.GeoKafkaMessage) error {
	payload, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", msg.UserId)),
		Value: payload,
	})
}
