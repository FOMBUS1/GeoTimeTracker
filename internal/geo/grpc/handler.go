package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/FOMBUS1/GeoTimeTracker/config"
	pb "github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_api"
	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type YandexGeocodeResponse struct {
	Response struct {
		GeoObjectCollection struct {
			FeatureMember []struct {
				GeoObject struct {
					MetaDataProperty struct {
						GeocoderMetaData struct {
							Address struct {
								Formatted string `json:"formatted"`
							} `json:"Address"`
						} `json:"GeocoderMetaData"`
					} `json:"metaDataProperty"`
				} `json:"GeoObject"`
			} `json:"featureMember"`
		} `json:"GeoObjectCollection"`
	} `json:"response"`
}


type Handler struct {
	pb.UnimplementedGeoServiceServer
	redis *redis.Client
	kafkaWriter *kafka.Writer
	cfg   *config.Config
}

func NewHandler(r *redis.Client, kw *kafka.Writer, cfg *config.Config) *Handler {
	return &Handler{
		redis: r,
		kafkaWriter: kw,
		cfg:   cfg,
	}
}

func (h *Handler) GetLocation(ctx context.Context, req *pb.LocationRequest) (*pb.LocationResponse, error) {
	log.Printf("Запрос получен: RequestID=%v, UserID=%v, Lat=%v, Long=%v", req.RequestId, req.UserId, req.Latitude, req.Longtitude)

	cacheKey := fmt.Sprintf("geo:%.5f:%.5f", req.Latitude, req.Longtitude)
	cachedAddress, err := h.redis.Get(ctx, cacheKey).Result()
	
	var address string
	if err == nil {
		log.Printf("Взят из кэша: %s", cachedAddress)
		address = cachedAddress
	} else {
		apiKey := os.Getenv("YANDEX_API_KEY")
		url := fmt.Sprintf(h.cfg.API.Url, apiKey, req.Longtitude, req.Latitude)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("ошибка запроса к Яндекс.Картам: %v", err)
		}
		defer resp.Body.Close()

		var geoData YandexGeocodeResponse
		if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
			return nil, fmt.Errorf("ошибка парсинга ответа Яндекса: %v", err)
		}

		address = "Адрес не найден"
		if len(geoData.Response.GeoObjectCollection.FeatureMember) > 0 {
			address = geoData.Response.GeoObjectCollection.FeatureMember[0].GeoObject.MetaDataProperty.GeocoderMetaData.Address.Formatted
		}

		h.redis.Set(ctx, cacheKey, address, 24*time.Hour)
		log.Printf("Определен через API и сохранен в кэш: %s", address)
	}

	geoEvent := &models.GeoModel{
		UserId:          req.UserId, 
		Departure:       req.Departure,
		LocationAddress: address,
		Time:            timestamppb.Now(),
	}
	if req.CustomLocation != nil {
		geoEvent.Location = req.CustomLocation
	}

	payload, err := proto.Marshal(geoEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal proto: %w", err)
	}

	err = h.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", geoEvent.UserId)),
		Value: payload,
	})

	if err != nil {
		log.Printf("Kafka Error: %v", err)
		return &pb.LocationResponse{
							Success: false,
						}, nil
	} else {
		log.Printf("Event sent to Kafka for User %d", geoEvent.UserId)
	}

	return &pb.LocationResponse{
		Success: true,
	}, nil
}