package main

import (
	"context"
	"fmt"
	"os"

	"github.com/FOMBUS1/GeoTimeTracker/config"
	bootstrap "github.com/FOMBUS1/GeoTimeTracker/internal/bootstrap/geo"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	ctx := context.Background()
	kw := bootstrap.NewKafkaWriter("user-geo-events", &cfg.Kafka)
	redisClient, err := bootstrap.NewRedisClient(ctx, &cfg.Redis)
	if err != nil {
		panic(fmt.Sprintf("ошибка загрузки redis, %v", err))
	}
	geoService := bootstrap.InitGeoService(redisClient, kw, cfg)
	geoServiceAPI := bootstrap.InitGeoServiceAPI(geoService)

	bootstrap.AppRun(*geoServiceAPI)
}
