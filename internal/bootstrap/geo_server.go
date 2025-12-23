package bootstrap

import (
	"github.com/FOMBUS1/GeoTimeTracker/config"
	"github.com/FOMBUS1/GeoTimeTracker/internal/clients/yandex"
	geoService "github.com/FOMBUS1/GeoTimeTracker/internal/services/geo"
	"github.com/FOMBUS1/GeoTimeTracker/internal/storage"
	kafka_storage "github.com/FOMBUS1/GeoTimeTracker/internal/storage/kafka"
	redis_storage "github.com/FOMBUS1/GeoTimeTracker/internal/storage/redis"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

func InitGeoService(rClient *redis.Client, kWriter *kafka.Writer, cfg *config.Config) *geoService.GeoService {
	cache := redis_storage.NewCache(rClient)
	producer := kafka_storage.NewProducer(kWriter)

	geocoder := yandex.NewClient(cfg.API.Key, cfg.API.Url)

	repo := storage.NewGeoRepository(cache, producer, geocoder)

	return geoService.NewGeoService(repo)
}
