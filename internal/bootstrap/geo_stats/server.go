package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	server "github.com/FOMBUS1/GeoTimeTracker/internal/api/geo_stats_api"
	usersgeoeventsconsumer "github.com/FOMBUS1/GeoTimeTracker/internal/consumer/users_geo_events_consumer"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_stats_api"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AppRun(api server.GeoStatsServiceAPI, usersGeoEventsConsumer *usersgeoeventsconsumer.UsersGeoEventsConsumer) {
	go usersGeoEventsConsumer.Consume(context.Background())
	go func() {
		if err := runGRPCServer(api); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %v", err))
		}
	}()

	if err := runGatewayServer(); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %v", err))
	}
}

func runGRPCServer(api server.GeoStatsServiceAPI) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	geo_stats_api.RegisterGeoStatsServiceServer(s, &api)

	slog.Info("gRPC-server server listening on :50051")
	return s.Serve(lis)
}

func runGatewayServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	swaggerPath := os.Getenv("swaggerPath")
	fmt.Println(swaggerPath)
	if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
		panic(fmt.Errorf("swagger file not found: %s", swaggerPath))
	}

	r := chi.NewRouter()
	r.Get("/swagger.json/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger.json", http.StatusMovedPermanently)
	})

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Swagger endpoint is triggered")
		http.ServeFile(w, r, swaggerPath)
	})

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusMovedPermanently)
	})

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := geo_stats_api.RegisterGeoStatsServiceHandlerFromEndpoint(ctx, mux, ":50051", opts)
	if err != nil {
		panic(err)
	}

	r.Mount("/", mux)

	slog.Info("gRPC-Gateway server listening on :8080")
	return http.ListenAndServe(":8080", r)
}
