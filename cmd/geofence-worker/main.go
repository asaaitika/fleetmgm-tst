package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/asaaitika/fleetmgm-tst/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		log.Fatalf("[GEOFENCE-WORKER][RABBITMQ][ERROR] >>> Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[GEOFENCE-WORKER][RABBITMQ][ERROR] >>> Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare queue (idempotent)
	q, err := ch.QueueDeclare(
		"geofence_alerts", // queue name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Fatalf("[GEOFENCE-WORKER][QUEUE][ERROR] >>> Failed to declare queue: %v", err)
	}

	// Set QoS (process 1 message at a time)
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("[GEOFENCE-WORKER][QOS][ERROR] >>> Failed to set QoS: %v", err)
	}

	// Register consumer
	msgs, err := ch.Consume(
		q.Name,                // queue
		"geofence_worker_001", // consumer tag
		false,                 // auto-ack (false = manual ack)
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		log.Fatalf("[GEOFENCE-WORKER][CONSUMER][ERROR] >>> Failed to register consumer: %v", err)
	}

	log.Println("[GEOFENCE-WORKER][RABBITMQ][INFO] >>> Connected to RabbitMQ")

	go func() {
		for msg := range msgs {
			processMessage(msg)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("[GEOFENCE-WORKER][APP][INFO] >>> Shutting down worker")
}

func processMessage(msg amqp.Delivery) {
	var event models.GeofenceEvent
	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		log.Printf("[GEOFENCE-WORKER][MSG][ERROR] >>> Failed to parse message: %v", err)
		msg.Nack(false, false)
		return
	}

	log.Printf("[GEOFENCE-WORKER][MSG][ALERT] >>> Vehicle %s entered %s at (%.4f, %.4f)",
		event.VehicleID, event.AreaName,
		event.Location.Latitude, event.Location.Longitude)

	// Simulate processing
	switch event.AreaName {
	case "Halte Bundaran HI":
		log.Printf("[GEOFENCE-WORKER][MSG][INFO] >>> SMS Alert: Vehicle %s at Halte Bundaran HI", event.VehicleID)
	case "Halte Blok M":
		log.Printf("[GEOFENCE-WORKER][MSG][INFO] >>> Sending notification: Vehicle %s arrived at Halte Blok M", event.VehicleID)
	case "Halte Cililitan (PGC)":
		log.Printf("[GEOFENCE-WORKER][MSG][INFO] >>> Dashboard updated: Vehicle %s reached Halte Cililitan (PGC)", event.VehicleID)
	case "Halte Cawang":
		log.Printf("[GEOFENCE-WORKER][MSG][INFO] >>> Email Alert: Vehicle %s at Halte Cawang", event.VehicleID)
	case "Halte Pulo Gadung":
		log.Printf("[GEOFENCE-WORKER][MSG][INFO] >>> Push notification: Vehicle %s reached Halte Pulo Gadung", event.VehicleID)
	}

	msg.Ack(false)
	log.Println("[GEOFENCE-WORKER][MSG][INFO] >>> Message processed and acknowledged")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
