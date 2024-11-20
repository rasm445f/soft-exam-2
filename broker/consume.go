package broker

import (
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

// Consume listens for messages on the specified queue and handles them using the provided handler.
func Consume(queue string, handler func(Event)) {
	channel := GetChannel()
	msgs, err := channel.Consume(
		queue,	// Queue Name
		"",		// Consumer Name
		true,	// Auto Acknowledge
		false,	// Exclusive
		false,	// No Local
		false,	// No Wait
		nil,	// Args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer key: %v", err)
	}

	go func() {
		for d:= range msgs {
			var event Event
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Failed to parse event: %v", err)
				continue
			}
			handler(event)
		}
	}()
}