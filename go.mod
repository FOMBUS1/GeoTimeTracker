module github.com/FOMBUS1/GeoTimeTracker

go 1.25.0

require (
	github.com/go-chi/chi/v5 v5.2.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	github.com/redis/go-redis/v9 v9.17.2
	github.com/samber/lo v1.52.0
	github.com/segmentio/kafka-go v0.4.49
	github.com/swaggo/http-swagger v1.3.4
	go.yaml.in/yaml/v4 v4.0.0-rc.3
	google.golang.org/genproto/googleapis/api v0.0.0-20250929231259-57b25ae835d4
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	github.com/swaggo/files v0.0.0-20220610200504-28940afbdbfe // indirect
	github.com/swaggo/swag v1.8.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250929231259-57b25ae835d4 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

tool (
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)
