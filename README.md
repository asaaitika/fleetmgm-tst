# 🚌 Fleet Tracking System

Sistem tracking armada bus real-time untuk monitoring dan management fleet dengan teknologi modern.

## 🏗️ Arsitektur Sistem

```
[GPS Device] → [MQTT Broker] → [Backend Services] → [PostgreSQL]
                     ↓                ↓
               [Mock Publisher]  [RabbitMQ] → [Geofence Worker]
                     ↓                ↓
               [Real-time Data]  [Alert System]
```

## 🛠️ Tech Stack

- **Backend**: Go (Gin Framework)
- **Database**: PostgreSQL 15
- **Message Brokers**: 
  - Eclipse Mosquitto (MQTT)
  - RabbitMQ (Event Processing)
- **Containerization**: Docker & Docker Compose
- **API**: RESTful API + Real-time MQTT

## 📋 Features Lengkap

- ✅ **Real-time GPS tracking** (interval 2 detik)
- ✅ **MQTT Integration** untuk IoT devices
- ✅ **Geofence detection & alerts** (terminal & halte)
- ✅ **Event-driven architecture** dengan RabbitMQ
- ✅ **Historical location data** storage
- ✅ **RESTful API endpoints** lengkap
- ✅ **Multi-service architecture** dengan Docker
- ✅ **Auto-scaling ready** dengan container orchestration

## 🚀 Quick Start

### Prerequisites

```bash
# Pastikan Docker & Docker Compose terinstall
docker --version          # Docker version 20.10+
docker-compose --version  # Docker Compose version 2.0+
```

### 1. Setup Project

```bash
# Clone dan masuk ke direktori project
cd fleet-tracking-system
```

### 2. Build & Start All Services

```bash
# Build dan start semua services sekaligus
docker-compose up --build

# Atau run di background
docker-compose up --build -d

# Monitor logs
docker-compose logs -f
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