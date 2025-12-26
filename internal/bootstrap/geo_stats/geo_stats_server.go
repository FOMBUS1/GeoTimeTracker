package bootstrap

import (
	"context"

	geoStatsService "github.com/FOMBUS1/GeoTimeTracker/internal/services/geoStats"
	"github.com/FOMBUS1/GeoTimeTracker/internal/storage/pgstorage"
)

func InitGeoStatsService(ctx context.Context, storage *pgstorage.PGstorage) *geoStatsService.GeoStatsService {
	return geoStatsService.NewGeoStatsService(ctx, storage)
}
