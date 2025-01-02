package broker

import (
	"fmt"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection
var channel *amqp091.Channel

func InitRabbitMQ() {
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	if rabbitMQHost == "" {
		rabbitMQHost = "localhost" // Default to localhost for local runs
	}
	var err error
	rabbitMQURL := fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitMQHost)

	conn, err = amqp091.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
}

func GetChannel() *amqp091.Channel {
	return channel
}

func CloseRabbitMQ() {
	if channel != nil {
		_ = channel.Close()
	}
	if conn != nil {
		_ = conn.Close()
	}
}

