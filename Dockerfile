# Base builder with dependencies
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod & go.sum for download dependencies
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build some services based on project structure
# CGO_ENABLED=0 for static binary (no need shared libraries)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-server cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mqtt-subscriber cmd/mqtt-subscriber/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mock-publisher cmd/mock-publisher/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o geofence-worker cmd/geofence-worker/main.go

# API Server
FROM alpine:latest AS api-server

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/api-server .

# Expose port
EXPOSE 8080

CMD ["./api-server"]

# MQTT Subscriber
FROM alpine:latest AS mqtt-subscriber

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/mqtt-subscriber .

CMD ["./mqtt-subscriber"]

# Mock Publisher  
FROM alpine:latest AS mock-publisher

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/mock-publisher .

CMD ["./mock-publisher"]

# Geofence Worker
FROM alpine:latest AS geofence-worker

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/geofence-worker .

CMD ["./geofence-worker"]