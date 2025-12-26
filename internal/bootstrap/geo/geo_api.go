package bootstrap

import (
	server "github.com/FOMBUS1/GeoTimeTracker/internal/api/geo_api"
	geoService "github.com/FOMBUS1/GeoTimeTracker/internal/services/geo"
)

func InitGeoServiceAPI(geoService *geoService.GeoService) *server.GeoServiceAPI {
	return server.NewGeoServiceAPI(geoService)
}
