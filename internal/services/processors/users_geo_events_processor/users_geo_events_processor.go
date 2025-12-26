package usersgeoeventsprocessor

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
)

type geoStatsService interface {
	UpsertUserVisits(ctx context.Context, usersGeoEvents []*models.GeoKafkaMessage) error
}

type UsersGeoEventsProcessor struct {
	geoStatsService geoStatsService
}

func NewUsersGeoEventsProcessor(geoStatsService geoStatsService) *UsersGeoEventsProcessor {
	return &UsersGeoEventsProcessor{
		geoStatsService: geoStatsService,
	}
}
