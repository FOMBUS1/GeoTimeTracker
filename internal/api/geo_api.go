package geo_api

import (
	"context"

	geo_models "github.com/FOMBUS1/GeoTimeTracker/internal/models"
	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_api"
)

type geoService interface {
	ProcessLocations(ctx context.Context, location []*geo_models.GeoModel) ([]bool, error)
}

type GeoServiceAPI struct {
	geo_api.UnimplementedGeoServiceServer
	geoService geoService
}

func NewGeoServiceAPI(geoService geoService) *GeoServiceAPI {
	return &GeoServiceAPI{
		geoService: geoService,
	}
}
