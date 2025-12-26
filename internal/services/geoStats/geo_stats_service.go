package geoStats

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_stats_api"
	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
)

type GeoStatsStorage interface {
	GetStatsByLocations(ctx context.Context, reqs []*geo_stats_api.UserLocationRequest) ([]*geo_stats_api.TimeSpentResponse, error)
	GetTopStats(ctx context.Context, reqs []*geo_stats_api.TimePeriodRequest) ([]*geo_stats_api.TimeSpentResponse, error)
	UpsertUserVisits(ctx context.Context, userVisits []*models.GeoKafkaMessage) error
}

type GeoStatsService struct {
	geoStatsStorage GeoStatsStorage
}

func NewGeoStatsService(ctx context.Context, geoStatsStorage GeoStatsStorage) *GeoStatsService {
	return &GeoStatsService{
		geoStatsStorage: geoStatsStorage,
	}
}
