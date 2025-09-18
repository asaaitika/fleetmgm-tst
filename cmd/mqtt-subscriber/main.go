package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asaaitika/fleetmgm-tst/internal/config"
	"github.com/asaaitika/fleetmgm-tst/internal/repositories"
	"github.com/asaaitika/fleetmgm-tst/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal("[MQTT-SUBCRIBER][DB][ERROR] >>> Failed to connect database:", err)
	}
	defer db.Close()

	log.Println("[MQTT-SUBCRIBER][DB][INFO] >>> Connected to PostgreSQL")

	vehicleRepo := repositories.NewVehicleRepository(db)

	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	var rabbitmqService *services.RabbitMQService
	for i := 0; i < 5; i++ {
		rabbitmqService, err = services.NewRabbitMQService(rabbitmqURL)
		if err == nil {
			break
		}
		log.Printf("[MQTT-SUBCRIBER][NEW-RABBITMQ][WARN] >>> RabbitMQ connection attempt %d failed: %v", i+1, err)
		if i < 4 {
			log.Println("[MQTT-SUBCRIBER][NEW-RABBITMQ][INFO] >>> Retrying in 5 seconds")
			time.Sleep(5 * time.Second)
		}
	}

	if rabbitmqService == nil {
		log.Println("[MQTT-SUBCRIBER][NEW-RABBITMQ][WARN] >>> Continuing without RabbitMQ - geofence events will not be published")
	} else {
		defer rabbitmqService.Close()
	}

	// Initialize MQTT service
	mqttBroker := getEnv("MQTT_BROKER", "tcp://localhost:1883")
	mqttService, err := services.NewMQTTService(mqttBroker, vehicleRepo)
	if err != nil {
		log.Fatal("[MQTT-SUBCRIBER][NEW-BROKER][ERROR] >>> Failed to create MQTT service:", err)
	}
	defer mqttService.Disconnect()

	// Initialize geofence service
	geoService := services.NewGeofenceService()
	mqttService.SetGeofenceService(geoService)

	if rabbitmqService != nil {
		mqttService.SetRabbitMQService(rabbitmqService)
		log.Println("[MQTT-SUBCRIBER][SET-RABBITMQ][INFO] >>> RabbitMQ service attached to MQTT service")
	} else {
		log.Println("[MQTT-SUBCRIBER][SET-RABBITMQ][WARN] >>> RabbitMQ service NOT attached - events won't be published")
	}

	// Subscribe to topics
	if err := mqttService.Subscribe(); err != nil {
		log.Fatal("[MQTT-SUBCRIBER][SUBSCRIBE][ERROR] >>> Failed to subscribe:", err)
	}

	log.Println("[MQTT-SUBCRIBER][APP][INFO] >>> MQTT Subscriber is running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("[MQTT-SUBCRIBER][APP][INFO] >>> Shutting down")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
