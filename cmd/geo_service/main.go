package main

import (
	"fmt"
	"log"
	"os"

	"github.com/FOMBUS1/GeoTimeTracker/config"
	"github.com/FOMBUS1/GeoTimeTracker/internal/geo/grpc"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

func main() {
	configPath := os.Getenv("configPath")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("ошибка парсинга конфига: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
    })

	kw := &kafka.Writer{
        Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
        Topic:    "user-geo-events",
        Balancer: &kafka.LeastBytes{},
    }
	defer kw.Close()

	handler := grpc.NewHandler(rdb, kw, cfg)
	
	grpcAddr := fmt.Sprintf(":%d", cfg.Server.GRPCPort)
    httpAddr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	
	srv := grpc.NewServer(grpcAddr, httpAddr, handler)

	log.Printf("Starting gRPC on %s and HTTP Gateway on %s...", grpcAddr, httpAddr)
	
	if err := srv.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}