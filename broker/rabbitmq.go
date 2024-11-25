package broker

import (
	"log"
	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection
var channel *amqp091.Channel

func InitRabbitMQ() {
	var err error
	conn, err = amqp091.Dial("amqp://guest:guest@localhost:5672/")
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