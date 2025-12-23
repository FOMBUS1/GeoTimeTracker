package geoService

import (
	"context"

	"github.com/FOMBUS1/GeoTimeTracker/internal/models"
	pb_models "github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GeoService) ProcessLocations(ctx context.Context, locations []*models.GeoModel) ([]bool, error) {
	results := make([]bool, len(locations))

	for i, loc := range locations {
		address, err := s.geoStorage.GetCachedAddress(ctx, loc.Latitude, loc.Longitude)
		if err != nil {
			address, err = s.geoStorage.FetchAddress(ctx, loc.Latitude, loc.Longitude)
			if err == nil {
				_ = s.geoStorage.SetCachedAddress(ctx, loc.Latitude, loc.Longitude, address)
			}
		}

		kafka_message := pb_models.GeoKafkaMessage{
			UserId:          loc.UserID,
			Departure:       loc.Departure,
			LocationAddress: address,
			Location:        &loc.CustomLocation,
			Time:            timestamppb.Now(),
		}

		err = s.geoStorage.SendToKafka(ctx, &kafka_message)

		results[i] = (err == nil)
	}

	return results, nil
}
