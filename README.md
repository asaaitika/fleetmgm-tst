# 🚌 Fleet Tracking System

Real-time fleet tracking system untuk armada bus menggunakan MQTT, PostgreSQL, dan RabbitMQ.

## 🚀 Tech Stack

- **Backend**: Go (Golang) 1.25
- **Database**: PostgreSQL 15
- **Message Queue**: RabbitMQ 3.12
- **MQTT Broker**: Eclipse Mosquitto 2.0
- **API Framework**: Gin
- **Container**: Docker & Docker Compose

## 🏗️ Arsitektur Sistem

```
[GPS Device] → [MQTT Broker] → [Subscriber Service] → [PostgreSQL]
↓
[Geofence Check]
↓
[RabbitMQ] → [Alert Worker]
```

## 📋 Features

- ✅ **Real-time GPS tracking** (interval 2 detik)
- ✅ **MQTT Integration** untuk IoT devices
- ✅ **Geofence detection & alerts** (terminal & halte)
- ✅ **Event-driven architecture** dengan RabbitMQ
- ✅ **Historical location data** storage
- ✅ **RESTful API endpoints** lengkap
- ✅ **Multi-service architecture** dengan Docker

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.25+ (untuk development)
- Make (optional)

### Installation & Running

1. **Clone repository**
```bash
git clone https://github.com/asaaitika/fleetmgm-tst.git
cd fleetmgm-tst
```

2. **Start all services**
```bash
# Build & start all services
docker-compose up --build

# or run di background
docker-compose up --build -d
```
3. **Run mock GPS publisher (testing)**
```bash
docker-compose --profile testing up -d
```

## 🔧 Development Mode

### Run Services Individually

```bash
# Start only infrastructure services
docker-compose up postgres rabbitmq mosquitto

# Run backend API locally
go run cmd/api/main.go

# Run MQTT subscriber locally
go run cmd/mqtt-subscriber/main.go

# Run mock publisher locally
go run cmd/mock-publisher/main.go

# Run geofence worker locally
go run cmd/geofence-worker/main.go
```

---

## 📄 License

Project ini dibuat untuk technical test - Fleet Tracking System.

**Author**: Permata Asa
**Version**: 1.0.0  
**Created**: September 2024

---