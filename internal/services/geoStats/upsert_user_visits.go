package geoStats

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
)

func (s *GeoStatsService) UpsertUserVisits(ctx context.Context, userVisits []*models.GeoKafkaMessage) error {
	return s.geoStatsStorage.UpsertUserVisits(ctx, userVisits)
}
