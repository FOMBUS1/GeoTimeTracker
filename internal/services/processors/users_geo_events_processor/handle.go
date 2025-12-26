package usersgeoeventsprocessor

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
)

func (p *UsersGeoEventsProcessor) Handle(ctx context.Context, usersGeoEvents *models.GeoKafkaMessage) error {
	return p.geoStatsService.UpsertUserVisits(ctx, []*models.GeoKafkaMessage{usersGeoEvents})
}
