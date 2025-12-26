package usersgeoeventsconsumer

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
)

type usersGeoEventsProcessor interface {
	Handle(ctx context.Context, studentsInfo *models.GeoKafkaMessage) error
}

type UsersGeoEventsConsumer struct {
	usersGeoEventsProcessor usersGeoEventsProcessor
	kafkaBroker             []string
	topicName               string
}

func NewUsersGeoEventsProcessor(usersGeoEventsProcessor usersGeoEventsProcessor, kafkaBroker []string, topicName string) *UsersGeoEventsConsumer {
	return &UsersGeoEventsConsumer{
		usersGeoEventsProcessor: usersGeoEventsProcessor,
		kafkaBroker:             kafkaBroker,
		topicName:               topicName,
	}
}
