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

✅ **Pelacakan GPS Real-time**: Menerima data lokasi setiap 2 detik.
✅ **Integrasi MQTT**: Dirancang untuk perangkat IoT (GPS Tracker).
✅ **Deteksi Geofence**: Memberikan notifikasi saat bus masuk ke area penting (terminal & halte).
✅ **Arsitektur Berbasis Event**: Menggunakan RabbitMQ untuk proses yang andal dan skalabel.
✅ **Penyimpanan Histori Lokasi**: Menyimpan jejak perjalanan untuk analisis.
✅ **RESTful API**: Menyediakan endpoint untuk mengakses data lokasi.
✅ **Arsitektur Multi-Service**: Setiap bagian sistem berjalan di kontainer Docker terpisah.

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.25+ (untuk development)
- Make (optional)

### 1. Installation & Running

```bash
# 1. Clone repository
git clone https://github.com/asaaitika/fleetmgm-tst.git
cd fleetmgm-tst

# 2. Build & Start all services
# Build & start all services
docker-compose up --build

# or run di background
docker-compose up --build -d
```

### 2. Migrasi Database (Manual)

```bash
# 1. Salin file migrasi SQL ke dalam kontainer postgres
docker cp internal/database/migrations.sql fleet_postgres:/tmp/migrations.sql

# 2. Eksekusi file migrasi tersebut
docker exec -it fleet_postgres psql -U fleet_admin -d fleet_db -f /tmp/migrations.sql

# 3. (Opsional) Verifikasi bahwa tabel sudah berhasil dibuat
docker exec -it fleet_postgres psql -U fleet_admin -d fleet_db -c "\dt"
```

### 3. Audit Log

```bash
# Melihat log dari semua service secara real-time
docker-compose logs -f

# Melihat log dari service tertentu (contoh: mqtt-subscriber)
docker-compose logs -f mqtt-subscriber
```

## 🧪 Testing with Mock Publisher

```bash
# Menjalankan semua service, TERMASUK mock-publisher
docker-compose --profile testing up --build -d

# Anda dapat melihat log dari publisher untuk memastikan data terkirim
docker-compose logs -f mock-publisher
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

## 📋 API Endpoints

- GET /vehicles/{vehicle_id}/location

Mengambil data lokasi terakhir dari kendaraan berdasarkan ID.

- GET /vehicles/{vehicle_id}/history?start=<timestamp>&end=<timestamp>

Mengambil histori perjalanan kendaraan dalam rentang waktu tertentu.

---

## 📄 License

Project ini dibuat untuk technical test - Fleet Tracking System.

**Author**: Permata Asa
**Version**: 1.0.0  
**Created**: September 2024

---