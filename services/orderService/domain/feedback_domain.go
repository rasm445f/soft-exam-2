package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/rasm445f/soft-exam-2/db/generated"
)

type FeedbackDomain struct {
	repo *generated.Queries
}

// NewFeedbackDomain initializes the domain layer
func NewFeedbackDomain(repo *generated.Queries) *FeedbackDomain {
	return &FeedbackDomain{repo: repo}
}

func (d *FeedbackDomain) GetAllFeedbacksDomain(ctx context.Context) ([]generated.Feedback, error) {
	feedbacks, err := d.repo.GetAllFeedbacks(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch feedbacks")
	}

	return feedbacks, nil
}

func (d *FeedbackDomain) GetFeedbackByOrderIdDomain(ctx context.Context, orderId int32) (*generated.Feedback, error) {
	feedback, err := d.repo.GetFeedbackByOrderId(ctx, orderId)
	if err != nil {
		return nil, errors.New("failed to get feedback by order")
	}

	return &feedback, nil
}

func (d *FeedbackDomain) CreateFeedbackDomain(ctx context.Context, feedbackParams generated.CreateFeedbackParams) (int32, error) {
	feedbackid, err := d.repo.CreateFeedback(ctx, feedbackParams)
	if err != nil {
		return 0, errors.New("failed to create feedback: " + err.Error())
	}

	err = d.UpdateDeliveryAgentRatingDomain(ctx, feedbackParams.Orderid)
	if err != nil {
		return 0, errors.New("failed to update average rating: " + err.Error())
	}

	return feedbackid, nil
}

func (d *FeedbackDomain) UpdateDeliveryAgentRatingDomain(ctx context.Context, orderId int32) error {
	// Fetch all feedbacks for all the delivery agent
	feedbacks, err := d.repo.GetAllFeedbacksFromDeliveryAgentByOrderId(ctx, orderId)
	if err != nil {
		return fmt.Errorf("failed to fetch feedbacks: %w", err)
	}

	// Calculate the average rating
	var totalRating float64
	var count int
	for _, feedback := range feedbacks {
		if feedback.Deliveryagentrating != nil {
			totalRating += float64(*feedback.Deliveryagentrating)
			count++
		}
	}

	if count == 0 {
		return fmt.Errorf("no ratings found for delivery agent %d", orderId)
	}

	avgRating := totalRating / float64(count)

	fmt.Print(avgRating)

	UpdateParams := generated.UpdateDeliveryAgentRatingParams{
		Rating: &avgRating,
		ID:     *feedbacks[0].Deliveryagentid,
	}

	fmt.Println(UpdateParams)

	// Update the delivery agent's rating
	err = d.repo.UpdateDeliveryAgentRating(ctx, UpdateParams)
	if err != nil {
		return fmt.Errorf("failed to update delivery agent rating: %w", err)
	}

	return nil
}

