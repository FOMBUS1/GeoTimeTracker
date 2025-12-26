package geoService

import (
	"context"
	"log/slog"
	"fmt"

	"github.com/FOMBUS1/GeoTimeTracker/internal/models"
	pb_models "github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GeoService) ProcessLocations(ctx context.Context, locations []*models.GeoModel) ([]bool, error) {
	results := make([]bool, len(locations))

	for i, loc := range locations {
		slog.Info("Checking if there address is already cached")
		address, err := s.geoStorage.GetCachedAddress(ctx, loc.Latitude, loc.Longitude)
		if err != nil {
			slog.Info("Address is not cached")
			slog.Info("Caching address")
			address, err = s.geoStorage.FetchAddress(ctx, loc.Latitude, loc.Longitude)
			if err == nil {
				slog.Info(fmt.Sprintf("Fetched address: %s", address))
				_ = s.geoStorage.SetCachedAddress(ctx, loc.Latitude, loc.Longitude, address)
			} else {
				slog.Error("Failed to fetch address from API", "err", err)
			}
			slog.Info("Address is cached")
		}

		slog.Info(fmt.Sprintf("Found address: %s", address))
		
		slog.Info("Sending message to kafka")
		kafka_message := pb_models.GeoKafkaMessage{
			UserId:          loc.UserID,
			Departure:       loc.Departure,
			LocationAddress: address,
			Location:        &loc.CustomLocation,
			Time:            timestamppb.Now(),
		}

		err = s.geoStorage.SendToKafka(ctx, &kafka_message)
		if err != nil {
			slog.Error("Failed to send message to Kafka", "err", err)
		} else {
			slog.Info("Message sent to kafka")
		}
		results[i] = (err == nil)
	}

	return results, nil
}
