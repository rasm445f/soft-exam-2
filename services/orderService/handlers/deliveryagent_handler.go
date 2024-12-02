package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/rasm445f/soft-exam-2/db/generated"
	"github.com/rasm445f/soft-exam-2/domain"
)

type DeliveryAgentHandler struct {
	domain *domain.DeliveryAgentDomain
}

func NewDeliveryAgentHandler(domain *domain.DeliveryAgentDomain) *DeliveryAgentHandler {
	return &DeliveryAgentHandler{domain: domain}
}

// GetAllDeliveryAgents godoc
//
// @Summary Get all deliveryAgents
// @Description Fetches a list of all Delivery Agents from the database
// @Tags DeliveryAgent CRUD
// @Produce application/json
// @Success 200 {array} generated.Deliveryagent
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/delivery-agent [get]
func (h *DeliveryAgentHandler) GetAllDeliveryAgents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		deliveryAgent, err := h.domain.GetAllDeliveryAgentsDomain(ctx)
		if err != nil {
			http.Error(w, "Failed to get deliveryAgent", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(deliveryAgent)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// GetDeliveryAgentById godoc
//
// @Summary Get deliveryAgent by deliveryAgent id
// @Description Fetches a deliveryAgent based on the id from the database
// @Tags DeliveryAgent CRUD
// @Produce application/json
// @Param deliveryAgentId path string true "DeliveryAgent ID"
// @Success 200 {object} generated.Deliveryagent
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/delivery-agent/{deliveryAgentId} [get]
func (h *DeliveryAgentHandler) GetDeliveryAgentById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		deliveryAgentIdStr := r.PathValue("deliveryAgentId")
		if deliveryAgentIdStr == "" {
			http.Error(w, "Missing DeliveryAgent Id path parameter", http.StatusBadRequest)
			return
		}

		deliveryAgentId, err := strconv.Atoi(deliveryAgentIdStr)
		if err != nil {
			http.Error(w, "Invalid DeliveryAgent ID", http.StatusBadRequest)
			return
		}

		feedback, err := h.domain.GetDeliveryAgentByIdDomain(ctx, int32(deliveryAgentId))
		if err != nil {
			http.Error(w, "DeliveryAgent not found", http.StatusNotFound)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(feedback)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// CreateDeliveryAgent godoc
//
// @Summary Create a new deliveryAgent
// @Description Creates a new deliveryAgent entry in the database
// @Tags DeliveryAgent CRUD
// @Accept  application/json
// @Produce application/json
// @Param deliveryAgent body generated.CreateDeliveryAgentParams true "DeliveryAgent object"
// @Success 201 {object} generated.Deliveryagent
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/delivery-agent [post]
func (h *DeliveryAgentHandler) CreateDeliveryAgent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var deliveryAgentParams generated.CreateDeliveryAgentParams

		err := json.NewDecoder(r.Body).Decode(&deliveryAgentParams)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Println(err)
			return
		}

		_, err = h.domain.CreateDeliveryAgentDomain(ctx, deliveryAgentParams)
		if err != nil {
			http.Error(w, "Failed to create deliveryAgent", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}
	
