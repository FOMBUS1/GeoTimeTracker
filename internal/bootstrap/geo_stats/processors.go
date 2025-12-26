package bootstrap

import (
	geoStatsService "github.com/FOMBUS1/GeoTimeTracker/internal/services/geoStats"
	usersgeoeventsprocessor "github.com/FOMBUS1/GeoTimeTracker/internal/services/processors/users_geo_events_processor"
)

func InitUsersGeoEventsProcessor(geoStatsService *geoStatsService.GeoStatsService) *usersgeoeventsprocessor.UsersGeoEventsProcessor {
	return usersgeoeventsprocessor.NewUsersGeoEventsProcessor(geoStatsService)
}
