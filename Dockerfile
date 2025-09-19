# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all services
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api-service ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o mqtt-subscriber ./cmd/mqtt-subscriber
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o geofence-worker ./cmd/geofence-worker
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o mock-publisher ./cmd/mock-publisher

# API Service
FROM alpine:3.18 AS api
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/api-service .
EXPOSE 8080
CMD ["./api-service"]

# MQTT Subscriber
FROM alpine:3.18 AS mqtt-subscriber
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/mqtt-subscriber .
CMD ["./mqtt-subscriber"]

# Geofence Worker
FROM alpine:3.18 AS geofence-worker
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/geofence-worker .
CMD ["./geofence-worker"]

# Mock Publisher
FROM alpine:3.18 AS mock-publisher
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/mock-publisher .
CMD ["./mock-publisher"]