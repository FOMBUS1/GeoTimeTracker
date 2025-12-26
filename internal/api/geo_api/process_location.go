package geo_api

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/models"
	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_api"
	"github.com/samber/lo"
)

func (s *GeoServiceAPI) GetLocation(ctx context.Context, req *geo_api.LocationRequests) (*geo_api.LocationResponse, error) {

	internalModels := mapLocations(req.Locations)

	results, err := s.geoService.ProcessLocations(ctx, internalModels)
	if err != nil {
		return nil, err
	}

	return &geo_api.LocationResponse{
		Success: results,
	}, nil
}

func mapLocations(pbItems []*geo_api.LocationRequest) []*models.GeoModel {
	return lo.Map(pbItems, func(item *geo_api.LocationRequest, _ int) *models.GeoModel {
		return &models.GeoModel{
			RequestID:      item.RequestId,
			UserID:         item.UserId,
			Longitude:      item.Longitude,
			Latitude:       item.Latitude,
			Departure:      item.Departure,
			CustomLocation: item.GetCustomLocation(),
		}
	})
}
