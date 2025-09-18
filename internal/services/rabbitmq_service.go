package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/asaaitika/fleetmgm-tst/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQService(url string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		"fleet.events", // exchange name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Declare queue
	_, err = ch.QueueDeclare(
		"geofence_alerts", // queue name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		"geofence_alerts", // queue name
		"geofence.#",      // routing key pattern
		"fleet.events",    // exchange
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	log.Println("[RABBITMQ-SERVICE][INFO] >>> Connected to RabbitMQ")

	return &RabbitMQService{
		conn:    conn,
		channel: ch,
	}, nil
}

// PublishEvent sends geofence event to RabbitMQ
func (s *RabbitMQService) PublishEvent(event *models.GeofenceEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	// Publish message
	err = s.channel.Publish(
		"fleet.events",   // exchange
		"geofence.entry", // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	log.Printf("[RABBITMQ-SERVICE][INFO] >>> Published event: Vehicle %s entered %s",
		event.VehicleID, event.AreaName)

	return nil
}

// Close closes RabbitMQ connection
func (s *RabbitMQService) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
	log.Println("[RABBITMQ-SERVICE][INFO] >>> Disconnected from RabbitMQ")
}
