package geoService

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
)

type GeoStorage interface {
	GetCachedAddress(ctx context.Context, lat, long float32) (string, error)
	SetCachedAddress(ctx context.Context, lat, long float32, address string) error
	FetchAddress(ctx context.Context, lat, long float32) (string, error)
	SendToKafka(ctx context.Context, event *models.GeoKafkaMessage) error
}

type GeoService struct {
	geoStorage GeoStorage
}

func NewGeoService(geoStorage GeoStorage) *GeoService {
	return &GeoService{
		geoStorage: geoStorage,
	}
}
