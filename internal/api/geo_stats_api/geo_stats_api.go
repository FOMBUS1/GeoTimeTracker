package geostatsapi

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_stats_api"
)

type geoStatsService interface {
	GetLocationStats(ctx context.Context, req *geo_stats_api.UserLocationRequests) (*geo_stats_api.TimeSpentResponses, error)
	GetTopTimeSpent(ctx context.Context, req *geo_stats_api.TimePeriodRequests) (*geo_stats_api.TimeSpentResponses, error)
}

type GeoStatsServiceAPI struct {
	geo_stats_api.UnimplementedGeoStatsServiceServer
	geoStatsService geoStatsService
}

func NewGeoStatsServiceAPI(geoStatsService geoStatsService) *GeoStatsServiceAPI {
	return &GeoStatsServiceAPI{
		geoStatsService: geoStatsService,
	}
}
