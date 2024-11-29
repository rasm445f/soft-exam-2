package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/rasm445f/soft-exam-2/db/generated"
	"github.com/rasm445f/soft-exam-2/domain"
)

type FeedbackHandler struct {
	domain *domain.FeedbackDomain
}

func NewFeedbackHandler(domain *domain.FeedbackDomain) *FeedbackHandler {
	return &FeedbackHandler{domain: domain}
}

// GetAllFeedbacks godoc
//
// @Summary Get all feedbacks
// @Description Fetches a list of all feedbacks from the database
// @Tags feedbacks
// @Produce application/json
// @Success 200 {array} generated.Feedback
// @Router /api/feedbacks [get]
func (h *FeedbackHandler) GetAllFeedbacks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		feedbacks, err := h.domain.GetAllFeedbacksDomain(ctx)
		if err != nil {
			http.Error(w, "Failed to get feedbacks", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(feedbacks)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// GetFeedbackByOrderId godoc
//
// @Summary Get feedback by order id
// @Description Fetches a feedback based on the order id from the database
// @Tags feedbacks
// @Produce application/json
// @Param orderId path string true "Order ID"
// @Success 200 {object} generated.Feedback
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Feedback not found"
// @Router /api/feedbacks/{orderId} [get]
func (h *FeedbackHandler) GetFeedbackByOrderId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orderIdStr := r.PathValue("orderId")
		if orderIdStr == "" {
			http.Error(w, "Missing Order Id path parameter", http.StatusBadRequest)
			return
		}

		orderId, err := strconv.Atoi(orderIdStr)
		if err != nil {
			http.Error(w, "Invalid Order ID", http.StatusBadRequest)
			return
		}

		feedback, err := h.domain.GetFeedbackByOrderIdDomain(ctx, int32(orderId))
		if err != nil {
			http.Error(w, "Feedback not found", http.StatusNotFound)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(feedback)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// // CreateFeedback godoc
// //
// // @Summary Create a new feedback
// // @Description Creates a new feedback entry in the database
// // @Tags feedbacks
// // @Accept  application/json
// // @Produce application/json
// // @Param feedback body generated.CreateFeedbackParams true "Feedback object"
// // @Success 201 {object} generated.Feedback
// // @Failure 400 {string} string "Bad request"
// // @Failure 500 {string} string "Internal server error"
// // @Router /api/feedback [post]
// func (h *FeedbackHandler) CreateFeedback() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()

// 		var feedbackParams generated.CreateFeedbackParams
// 		err := json.NewDecoder(r.Body).Decode(&feedbackParams)
// 		if err != nil {
// 			http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 			log.Println(err)
// 			return
// 		}

// 		_, err = h.domain.CreateFeedbackDomain(ctx, feedbackParams)
// 		if err != nil {
// 			http.Error(w, "Failed to create feedback", http.StatusInternalServerError)
// 			log.Println(err)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusCreated)
// 	}
// }
	

// CreateFeedback godoc
//
// @Summary Create a new feedback
// @Description Creates a new feedback entry in the database
// @Tags feedbacks
// @Accept  application/json
// @Produce application/json
// @Param feedback body generated.CreateFeedbackParams true "Feedback object"
// @Success 201 {object} generated.Feedback
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/feedback [post]
func (h *FeedbackHandler) CreateFeedback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var feedbackParams generated.CreateFeedbackParams
		err := json.NewDecoder(r.Body).Decode(&feedbackParams)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Println(err)
			return
		}

		_, err = h.domain.CreateFeedbackDomain(ctx, feedbackParams)
		if err != nil {
			http.Error(w, "Failed to create feedback", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}
