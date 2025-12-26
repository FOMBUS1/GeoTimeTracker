package main

import (
	"context"
	"fmt"
	"os"

	"github.com/FOMBUS1/GeoTimeTracker/config"
	bootstrap "github.com/FOMBUS1/GeoTimeTracker/internal/bootstrap/geo_stats"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	storage := bootstrap.InitPGStorage(cfg)
	geoStatsService := bootstrap.InitGeoStatsService(context.Background(), storage)
	geoStatsServiceAPI := bootstrap.InitGeoStatsServiceAPI(geoStatsService)

	usersGeoEventsProcessor := bootstrap.InitUsersGeoEventsProcessor(geoStatsService)
	usersGeoEventsConsumer := bootstrap.InitUsersGeoEventsConsumer(cfg, usersGeoEventsProcessor)

	bootstrap.AppRun(*geoStatsServiceAPI, usersGeoEventsConsumer)
}
