package grpc

import (
	"context"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	
	pb "github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_api"
)

type Server struct {
	grpcAddr string
	httpAddr string
	handler  *Handler
}

func NewServer(grpcAddr, httpAddr string, h *Handler) *Server {
	return &Server{grpcAddr: grpcAddr, 
				httpAddr: httpAddr, 
				handler:  h}
}

func (s *Server) Run() error {
	gServer := grpc.NewServer()
	pb.RegisterGeoServiceServer(gServer, s.handler)

	lis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return err
	}

	go func() {
		if err := gServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err = pb.RegisterGeoServiceHandlerFromEndpoint(ctx, mux, s.grpcAddr, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(s.httpAddr, mux)
}