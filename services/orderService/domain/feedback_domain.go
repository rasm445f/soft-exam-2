package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rasm445f/soft-exam-2/db/generated"
)

type FeedbackDomain struct {
	repo *generated.Queries
}

// NewOrderDomain initializes the domain layer
func NewFeedbackDomain(repo *generated.Queries) *FeedbackDomain {
	return &FeedbackDomain{repo: repo}
}

func (d *FeedbackDomain) CreateFeedbackDomain(ctx context.Context, feedbackParams generated.CreateFeedbackParams)