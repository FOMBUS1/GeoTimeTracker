package grpc

import (
	stdgrpc "google.golang.org/grpc"

	"github.com/FOMBUS1/GeoTimeTracker/api/geo_api"
)

type Server struct {
	geo_api.UnimplementedGeoServer
	server *stdgrpc.Server
}
