package bootstrap

import (
	"fmt"

	"github.com/FOMBUS1/GeoTimeTracker/config"
	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(topicName string, cfg *config.KafkaConfig) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
	}
}
