package geoStats

import (
	"errors"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_stats_api"
	"gotest.tools/v3/assert"
)

func (s *GeoStatsServiceSuite) TestGetTopTimeSpentSuccess() {
	req := &geo_stats_api.TimePeriodRequests{
		Requests: []*geo_stats_api.TimePeriodRequest{
			{
				UserId: 1,
				Period: geo_stats_api.Period_WEEK,
				TopK:   ptrUint32(5),
			},
		},
	}

	expectedStats := []*geo_stats_api.TimeSpentResponse{
		{
			UserId: 1,
			Stats: []*geo_stats_api.LocationStat{
				{
					LocationId:  101,
					TimeSeconds: 3600,
					Period:      geo_stats_api.Period_WEEK,
				},
			},
		},
	}

	s.geoStatsStorage.EXPECT().
		GetTopStats(s.ctx, req.Requests).
		Return(expectedStats, nil)

	resp, err := s.geoStatsService.GetTopTimeSpent(s.ctx, req)

	assert.NilError(s.T(), err)
	assert.Equal(s.T(), len(resp.Responses), len(expectedStats))
	assert.Equal(s.T(), resp.Responses[0].UserId, expectedStats[0].UserId)
	assert.Equal(s.T(), resp.Responses[0].Stats[0].LocationId, expectedStats[0].Stats[0].LocationId)
}

func (s *GeoStatsServiceSuite) TestGetTopTimeSpentStorageError() {
	req := &geo_stats_api.TimePeriodRequests{
		Requests: []*geo_stats_api.TimePeriodRequest{{UserId: 1}},
	}
	wantErr := errors.New("internal storage error")

	s.geoStatsStorage.EXPECT().
		GetTopStats(s.ctx, req.Requests).
		Return(nil, wantErr)

	resp, err := s.geoStatsService.GetTopTimeSpent(s.ctx, req)

	assert.ErrorIs(s.T(), err, wantErr)
	assert.Assert(s.T(), resp == nil)
}

func (s *GeoStatsServiceSuite) TestGetTopTimeSpentEmptyRequest() {
	req := &geo_stats_api.TimePeriodRequests{
		Requests: []*geo_stats_api.TimePeriodRequest{},
	}

	s.geoStatsStorage.EXPECT().
		GetTopStats(s.ctx, req.Requests).
		Return([]*geo_stats_api.TimeSpentResponse{}, nil)

	resp, err := s.geoStatsService.GetTopTimeSpent(s.ctx, req)

	assert.NilError(s.T(), err)
	assert.Equal(s.T(), len(resp.Responses), 0)
}

func ptrUint32(v uint32) *uint32 {
	return &v
}
