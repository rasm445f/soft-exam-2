package domain

import (
	"context"
	"errors"

	"github.com/rasm445f/soft-exam-2/db/generated"
)

type DeliveryAgentDomain struct {
	repo *generated.Queries
}

// NewDeliveryAgentDomain initializes the domain layer
func NewDeliveryAgentDomain(repo *generated.Queries) *DeliveryAgentDomain {
	return &DeliveryAgentDomain{repo: repo}
}

func (d *DeliveryAgentDomain) GetAllDeliveryAgentsDomain(ctx context.Context) ([]generated.Deliveryagent, error) {
	deliveryAgents, err := d.repo.GetAllDeliveryAgents(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch deliveryAgents")
	}

	return deliveryAgents, nil
}

func (d *DeliveryAgentDomain) GetDeliveryAgentByIdDomain(ctx context.Context, deliveryAgentId int32) (*generated.Deliveryagent, error) {
	deliveryAgent, err := d.repo.GetDeliveryAgentById(ctx, deliveryAgentId)
	if err != nil {
		return nil, errors.New("failed to get deliveryAgent by id")
	}

	return &deliveryAgent, nil
}

func (d *DeliveryAgentDomain) CreateDeliveryAgentDomain(ctx context.Context, deliveryAgentParams generated.CreateDeliveryAgentParams) (int32, error) {
	deliveryAgentId, err := d.repo.CreateDeliveryAgent(ctx, deliveryAgentParams)
	if err != nil {
		return 0, errors.New("failed to create delivery agent: " + err.Error())
	}
	return deliveryAgentId, nil
}

// func (d *DeliveryAgentDomain) UpdateDeliveryAgentRatingDomain(ctx context.Context, deliveryAgentId int32) error {
// 	// Fetch all feedbacks for all the delivery agent
// 	feedbacks, err := d.repo.GetAllFeedbacksByDeliveryAgentId(ctx, &deliveryAgentId)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch feedbacks: %w", err)
// 	}

// 	// Calculate the average rating
// 	var totalRating float64
// 	var count int
// 	for _, feedback := range feedbacks {
// 		if feedback.Deliveryagentrating != nil {
// 			totalRating += float64(*feedback.Deliveryagentrating)
// 			count++
// 		}
// 	}

// 	if count == 0 {
// 		return fmt.Errorf("no ratings found for delivery agent %d", deliveryAgentId)
// 	}

// 	avgRating := totalRating / float64(count)

// 	UpdateParams := generated.UpdateDeliveryAgentRatingParams {
// 		Rating: &avgRating,
// 		ID: deliveryAgentId,
// 	}

// 	// Update the delivery agent's rating
// 	err = d.repo.UpdateDeliveryAgentRating(ctx, UpdateParams)
// 	if err != nil {
// 		return fmt.Errorf("failed to update delivery agent rating: %w", err)
// 	}

// 	return nil
// }
