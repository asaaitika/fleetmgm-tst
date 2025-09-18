package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/asaaitika/fleetmgm-tst/internal/models"
	"github.com/asaaitika/fleetmgm-tst/internal/repositories"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTService struct {
	client   mqtt.Client
	repo     *repositories.VehicleRepository
	geofence *GeofenceService
	rabbitmq *RabbitMQService
}

func NewMQTTService(broker string, repo *repositories.VehicleRepository) (*MQTTService, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("fleet_subscriber_001")
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)

	opts.OnConnect = func(client mqtt.Client) {
		log.Println("[MQTT-SERVICE][INFO] >>> Connected to MQTT broker")
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("[MQTT-SERVICE][ERROR] >>> MQTT connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)

	// Connect to broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &MQTTService{
		client: client,
		repo:   repo,
	}, nil
}

// SetGeofenceService inject geofence service
func (s *MQTTService) SetGeofenceService(geo *GeofenceService) {
	s.geofence = geo
}

// SetRabbitMQService inject rabbitmq service
func (s *MQTTService) SetRabbitMQService(rmq *RabbitMQService) {
	s.rabbitmq = rmq
	if rmq != nil {
		log.Println("[MQTT-SERVICE][INFO] >>> RabbitMQ service successfully attached")
	} else {
		log.Println("[MQTT-SERVICE][WARN] >>> Nil RabbitMQ service provided")
	}
}

// Subscribe to vehicle location topic
func (s *MQTTService) Subscribe() error {
	topic := "/fleet/vehicle/+/location"

	// Subscribe with QoS 1 (at least once delivery)
	token := s.client.Subscribe(topic, 1, s.handleMessage)

	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	log.Printf("[MQTT-SERVICE][INFO] >>> Subscribed to topic: %s", topic)
	return nil
}

// handleMessage process incoming MQTT messages
func (s *MQTTService) handleMessage(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[MQTT-SERVICE][DEBUG] >>> Received message on topic: %s", msg.Topic())

	var payload models.MQTTPayload
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		log.Printf("[MQTT-SERVICE][ERROR] >>> Failed to parse payload: %v", err)
		return
	}

	if err := s.validatePayload(&payload); err != nil {
		log.Printf("[MQTT-SERVICE][ERROR] >>> Invalid payload: %v", err)
		return
	}

	if err := s.repo.InsertLocation(&payload); err != nil {
		log.Printf("[MQTT-SERVICE][ERROR] >>> Failed to save to database: %v", err)
		return
	}

	log.Printf("[MQTT-SERVICE][INFO] >>> Saved location for vehicle: %s", payload.VehicleID)

	if s.geofence != nil {
		s.checkGeofence(&payload)
	}
}

// validatePayload validates incoming data
func (s *MQTTService) validatePayload(payload *models.MQTTPayload) error {
	if payload.VehicleID == "" {
		return fmt.Errorf("vehicle_id is required")
	}

	// Validate koordinat Indonesia (approximate)
	if payload.Latitude < -11 || payload.Latitude > 6 {
		return fmt.Errorf("latitude out of range for Indonesia")
	}

	if payload.Longitude < 95 || payload.Longitude > 141 {
		return fmt.Errorf("longitude out of range for Indonesia")
	}

	if payload.Timestamp <= 0 {
		return fmt.Errorf("invalid timestamp")
	}

	return nil
}

// checkGeofence checks if vehicle entered geofence area
func (s *MQTTService) checkGeofence(payload *models.MQTTPayload) {
	areas, err := s.repo.GetGeofenceAreas()
	if err != nil {
		log.Printf("[MQTT-SERVICE][ERROR] >>> Failed to get geofence areas: %v", err)
		return
	}

	for _, area := range areas {
		distance := s.geofence.CalculateDistance(
			payload.Latitude, payload.Longitude,
			area.CenterLatitude, area.CenterLongitude,
		)

		// Check if within radius
		if distance <= float64(area.RadiusMeters) {
			log.Printf("[MQTT-SERVICE][ALERT] >>> Vehicle %s entered geofence: %s (%.2fm away)",
				payload.VehicleID, area.Name, distance)

			if s.rabbitmq != nil {
				event := models.GeofenceEvent{
					VehicleID: payload.VehicleID,
					Event:     "geofence_entry",
					Location: models.Location{
						Latitude:  payload.Latitude,
						Longitude: payload.Longitude,
					},
					Timestamp: payload.Timestamp,
					AreaName:  area.Name,
				}

				if err := s.rabbitmq.PublishEvent(&event); err != nil {
					log.Printf("[MQTT-SERVICE][ERROR] >>> Failed to publish event: %v", err)
				}
			} else {
				log.Printf("[MQTT-SERVICE][WARN] >>> RabbitMQ not connected, skipping event publish")
			}
		}
	}
}

// Disconnect from MQTT broker
func (s *MQTTService) Disconnect() {
	s.client.Disconnect(250)
	log.Println("[MQTT-SERVICE][INFO] >>> Disconnected from MQTT broker")
}
