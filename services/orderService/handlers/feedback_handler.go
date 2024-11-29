package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/rasm445f/soft-exam-2/broker"
	"github.com/rasm445f/soft-exam-2/db/generated"
	"github.com/rasm445f/soft-exam-2/domain"
)

type FeedbackHandler struct {
	domain *domain.FeedbackDomain
}

func NewFeedbackHandler(domain *domain.FeedbackDomain) *FeedbackHandler {
	return &FeedbackHandler{domain: domain}
}