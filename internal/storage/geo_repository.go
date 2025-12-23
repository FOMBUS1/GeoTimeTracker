package storage

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/clients/yandex"
	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
	"github.com/FOMBUS1/GeoTimeTracker/internal/storage/kafka"
	"github.com/FOMBUS1/GeoTimeTracker/internal/storage/redis"
)

type GeoRepository struct {
	cache    *redis.Cache
	producer *kafka.Producer
	geocoder *yandex.Client
}

func NewGeoRepository(c *redis.Cache, p *kafka.Producer, g *yandex.Client) *GeoRepository {
	return &GeoRepository{cache: c, producer: p, geocoder: g}
}

func (r *GeoRepository) GetCachedAddress(ctx context.Context, lat, long float32) (string, error) {
	return r.cache.Get(ctx, lat, long)
}

func (r *GeoRepository) SetCachedAddress(ctx context.Context, lat, long float32, addr string) error {
	return r.cache.Set(ctx, lat, long, addr)
}

func (r *GeoRepository) FetchAddress(ctx context.Context, lat, long float32) (string, error) {
	return r.geocoder.Fetch(ctx, lat, long)
}

func (r *GeoRepository) SendToKafka(ctx context.Context, event *models.GeoKafkaMessage) error {
	return r.producer.Send(ctx, event)
}
