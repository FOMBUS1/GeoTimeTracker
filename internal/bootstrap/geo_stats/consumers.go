package bootstrap

import (
	"fmt"

	"github.com/FOMBUS1/GeoTimeTracker/config"
	usersgeoeventsconsumer "github.com/FOMBUS1/GeoTimeTracker/internal/consumer/users_geo_events_consumer"
	usersgeoeventsprocessor "github.com/FOMBUS1/GeoTimeTracker/internal/services/processors/users_geo_events_processor"
)

func InitUsersGeoEventsConsumer(cfg *config.Config, userGeoEventsProcessor *usersgeoeventsprocessor.UsersGeoEventsProcessor) *usersgeoeventsconsumer.UsersGeoEventsConsumer {
	kafkaBrockers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}

	return usersgeoeventsconsumer.NewUsersGeoEventsProcessor(userGeoEventsProcessor, kafkaBrockers, "user-geo-events")
}
