package bootstrap

import (
	server "github.com/FOMBUS1/GeoTimeTracker/internal/api/geo_stats_api"
	geoStatsService "github.com/FOMBUS1/GeoTimeTracker/internal/services/geoStats"
)

func InitGeoStatsServiceAPI(geoStatsService *geoStatsService.GeoStatsService) *server.GeoStatsServiceAPI {
	return server.NewGeoStatsServiceAPI(geoStatsService)
}
