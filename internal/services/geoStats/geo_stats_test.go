package geoStats

import (
	"errors"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_stats_api"
	"gotest.tools/v3/assert"
)

func (s *GeoStatsServiceSuite) TestGetLocationStatsSuccess() {
	req := &geo_stats_api.UserLocationRequests{
		Requests: []*geo_stats_api.UserLocationRequest{
			{
				UserId:     1,
				LocationId: 101,
				Period:     geo_stats_api.Period_MONTH,
			},
		},
	}

	expectedStats := []*geo_stats_api.TimeSpentResponse{
		{
			UserId: 1,
			Stats: []*geo_stats_api.LocationStat{
				{
					LocationId:  101,
					TimeSeconds: 7200,
					Period:      geo_stats_api.Period_MONTH,
				},
			},
		},
	}

	s.geoStatsStorage.EXPECT().
		GetStatsByLocations(s.ctx, req.Requests).
		Return(expectedStats, nil)

	resp, err := s.geoStatsService.GetLocationStats(s.ctx, req)

	assert.NilError(s.T(), err)
	assert.Assert(s.T(), resp != nil)
	assert.Equal(s.T(), len(resp.Responses), 1)
	assert.Equal(s.T(), resp.Responses[0].UserId, uint64(1))
	assert.Equal(s.T(), resp.Responses[0].Stats[0].TimeSeconds, uint64(7200))
}

func (s *GeoStatsServiceSuite) TestGetLocationStatsStorageError() {
	req := &geo_stats_api.UserLocationRequests{
		Requests: []*geo_stats_api.UserLocationRequest{
			{UserId: 1, LocationId: 101},
		},
	}
	wantErr := errors.New("database connection error")

	s.geoStatsStorage.EXPECT().
		GetStatsByLocations(s.ctx, req.Requests).
		Return(nil, wantErr)

	resp, err := s.geoStatsService.GetLocationStats(s.ctx, req)

	assert.ErrorIs(s.T(), err, wantErr)
	assert.Assert(s.T(), resp == nil)
}

func (s *GeoStatsServiceSuite) TestGetLocationStatsMultipleRequests() {
	req := &geo_stats_api.UserLocationRequests{
		Requests: []*geo_stats_api.UserLocationRequest{
			{UserId: 1, LocationId: 101},
			{UserId: 2, LocationId: 202},
		},
	}

	s.geoStatsStorage.EXPECT().
		GetStatsByLocations(s.ctx, req.Requests).
		Return([]*geo_stats_api.TimeSpentResponse{
			{UserId: 1},
			{UserId: 2},
		}, nil)

	resp, err := s.geoStatsService.GetLocationStats(s.ctx, req)

	assert.NilError(s.T(), err)
	assert.Equal(s.T(), len(resp.Responses), 2)
}
