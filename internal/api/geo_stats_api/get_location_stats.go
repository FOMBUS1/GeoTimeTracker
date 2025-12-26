package geostatsapi

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_stats_api"
)

func (s *GeoStatsServiceAPI) GetLocationStats(ctx context.Context, req *geo_stats_api.UserLocationRequests) (*geo_stats_api.TimeSpentResponses, error) {
	results, err := s.geoStatsService.GetLocationStats(ctx, req)

	if err != nil {
		return nil, err
	}

	return results, nil
}
