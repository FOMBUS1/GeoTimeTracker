package geoStats

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_stats_api"
)

func (s *GeoStatsService) GetLocationStats(ctx context.Context, req *geo_stats_api.UserLocationRequests) (*geo_stats_api.TimeSpentResponses, error) {
	response, err := s.geoStatsStorage.GetStatsByLocations(ctx, req.Requests)
	if err != nil {
		return nil, err
	}

	return &geo_stats_api.TimeSpentResponses{
		Responses: response,
	}, err
}
