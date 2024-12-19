package domain

import "context"

// Broker defines an interface for publishing events
type Broker interface {
    Publish(ctx context.Context, queue string, event Event) error
}

// Event represents a message to be published
type Event struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}