FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o geo-service ./cmd/geo_service

ENV configPath=config.yaml

EXPOSE 8080 8082
CMD ["/app/geo-service"]