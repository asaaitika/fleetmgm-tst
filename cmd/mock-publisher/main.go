package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/asaaitika/fleetmgm-tst/internal/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Coordinate area Jakarta for simulation
var jakartaBounds = struct {
	MinLat, MaxLat float64
	MinLon, MaxLon float64
}{
	MinLat: -6.3000,  // Jakarta Selatan
	MaxLat: -6.1000,  // Jakarta Utara
	MinLon: 106.7000, // Jakarta Barat
	MaxLon: 106.9000, // Jakarta Timur
}

// Route simulation
var simulationRoute = []struct {
	Lat  float64
	Lon  float64
	Name string
}{
	{-6.1944, 106.8229, "Halte Bundaran HI"},
	{-6.2431, 106.8018, "Halte Blok M"},
	{-6.2460, 106.8990, "Halte Cililitan (PGC)"},
	{-6.2380, 106.8550, "Halte Cawang"},
	{-6.1920, 106.8950, "Halte Pulo Gadung"},
}

func main() {
	broker := getEnv("MQTT_BROKER", "tcp://localhost:1883")
	interval := getEnv("PUBLISH_INTERVAL", "2s")
	vehicleIDs := strings.Split(getEnv("VEHICLE_IDS", "B1234XYZ,B5678ABC"), ",")

	publishInterval, err := time.ParseDuration(interval)
	if err != nil {
		log.Fatal("[MOCK-PUBLISHER]][APP][ERROR] >>> Invalid interval:", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(fmt.Sprintf("mock_publisher_%d", time.Now().Unix()))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("[MOCK-PUBLISHER]][MQTT][ERROR] >>> Failed to connect:", token.Error())
	}
	defer client.Disconnect(250)

	log.Printf("[MOCK-PUBLISHER]][MQTT][INFO] >>> Connected to MQTT broker: %s", broker)
	log.Printf("[MOCK-PUBLISHER]][MQTT][INFO] >>> Publishing data every %s for vehicles: %v", interval, vehicleIDs)

	// Random seed
	rand.Seed(time.Now().UnixNano())

	// Track position for every vehicle
	vehiclePositions := make(map[string]int)
	for _, vid := range vehicleIDs {
		vehiclePositions[vid] = rand.Intn(len(simulationRoute))
	}

	ticker := time.NewTicker(publishInterval)
	defer ticker.Stop()

	log.Println("[MOCK-PUBLISHER]][INFO] >>> Starting to publish mock GPS data")

	for range ticker.C {
		for _, vehicleID := range vehicleIDs {
			// Get current position
			pos := vehiclePositions[vehicleID]
			currentPoint := simulationRoute[pos]

			// Small random offset
			lat := currentPoint.Lat + (rand.Float64()-0.5)*0.0001
			lon := currentPoint.Lon + (rand.Float64()-0.5)*0.0001

			payload := models.MQTTPayload{
				VehicleID: vehicleID,
				Latitude:  lat,
				Longitude: lon,
				Timestamp: time.Now().Unix(),
			}

			data, err := json.Marshal(payload)
			if err != nil {
				log.Printf("[MOCK-PUBLISHER]][ERROR] >>> Failed to marshal: %v", err)
				continue
			}

			// Publish to MQTT
			topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID)
			token := client.Publish(topic, 1, false, data)
			token.Wait()

			if token.Error() != nil {
				log.Printf("[MOCK-PUBLISHER]][ERROR] >>> Failed to publish: %v", token.Error())
			} else {
				log.Printf("[MOCK-PUBLISHER]][DEBUG] >>> Published: Vehicle=%s, Location=%s (%.4f,%.4f)",
					vehicleID, currentPoint.Name, lat, lon)
			}

			// Move to next position
			vehiclePositions[vehicleID] = (pos + 1) % len(simulationRoute)
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
