// package studentsService
package geoStats

import (
	"context"
	"errors"
	"testing"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
	"github.com/FOMBUS1/GeoTimeTracker/internal/services/geoStats/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gotest.tools/v3/assert"
)

type GeoStatsServiceSuite struct {
	suite.Suite
	ctx             context.Context
	geoStatsStorage *mocks.GeoStatsStorage
	geoStatsService *GeoStatsService
}

func (s *GeoStatsServiceSuite) SetupTest() {
	s.geoStatsStorage = mocks.NewGeoStatsStorage(s.T())
	s.ctx = context.Background()
	s.geoStatsService = NewGeoStatsService(s.ctx, s.geoStatsStorage)
}

func (s *GeoStatsServiceSuite) TestUpsertUserVisitsWithoutLocationSuccess() {
	data := []*models.GeoKafkaMessage{
		{
			UserId:          1,
			Departure:       false,
			LocationAddress: "Ставропольская улица, 13",
			Location:        nil,
			Time:            timestamppb.Now(),
		},
	}
	s.geoStatsStorage.EXPECT().UpsertUserVisits(s.ctx, data).Return(nil)

	err := s.geoStatsService.UpsertUserVisits(s.ctx, data)

	assert.NilError(s.T(), err)
}

func (s *GeoStatsServiceSuite) TestUpsertUserVisitsWithLocationSuccess() {
	location := "Дом"
	data := []*models.GeoKafkaMessage{
		{
			UserId:          1,
			Departure:       false,
			LocationAddress: "Ставропольская улица, 13",
			Location:        &location,
			Time:            timestamppb.Now(),
		},
	}
	s.geoStatsStorage.EXPECT().
		UpsertUserVisits(s.ctx, data).
		Return(nil)

	err := s.geoStatsService.UpsertUserVisits(s.ctx, data)

	assert.NilError(s.T(), err)
}

func (s *GeoStatsServiceSuite) TestUpsertUserVisitsStorageError() {
	data := []*models.GeoKafkaMessage{
		{UserId: 1, LocationAddress: "Улица Пушкина"},
	}
	wantErr := errors.New("db connection refused")

	s.geoStatsStorage.EXPECT().UpsertUserVisits(s.ctx, data).Return(wantErr)

	err := s.geoStatsService.UpsertUserVisits(s.ctx, data)

	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *GeoStatsServiceSuite) TestUpsertUserVisitsMultipleMessages() {
	data := []*models.GeoKafkaMessage{
		{UserId: 1, Departure: false, LocationAddress: "Ставропольская улица, 13"},
		{UserId: 2, Departure: true, LocationAddress: "Тюмень, улица Герцена, 72с1"},
	}

	s.geoStatsStorage.EXPECT().UpsertUserVisits(s.ctx, data).Return(nil)
	err := s.geoStatsService.UpsertUserVisits(s.ctx, data)

	assert.NilError(s.T(), err)
}

func (s *GeoStatsServiceSuite) TestUpsertUserVisitsEmptyData() {
	var data []*models.GeoKafkaMessage

	s.geoStatsStorage.EXPECT().
		UpsertUserVisits(s.ctx, data).
		Return(nil)

	err := s.geoStatsService.UpsertUserVisits(s.ctx, data)

	assert.NilError(s.T(), err)
}

func (s *GeoStatsServiceSuite) TestUpsertUserVisitsAnythingMatch() {
	s.geoStatsStorage.EXPECT().UpsertUserVisits(s.ctx, mock.Anything).Return(nil)

	err := s.geoStatsService.UpsertUserVisits(s.ctx, []*models.GeoKafkaMessage{{UserId: 99}})

	assert.NilError(s.T(), err)
}

func (s *GeoStatsServiceSuite) TestUpsertUserVisitsNilMessageInSlice() {
	data := []*models.GeoKafkaMessage{nil}

	s.geoStatsStorage.EXPECT().UpsertUserVisits(s.ctx, data).Return(errors.New("invalid data"))

	err := s.geoStatsService.UpsertUserVisits(s.ctx, data)

	assert.Check(s.T(), err != nil)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(GeoStatsServiceSuite))
}
