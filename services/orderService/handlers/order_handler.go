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

type OrderHandler struct {
	domain *domain.OrderDomain
}

func NewOrderHandler(domain *domain.OrderDomain) *OrderHandler {
	return &OrderHandler{domain: domain}
}

// GetAllOrders godoc
//
// @Summary Get all orders
// @Description Fetches a list of all orders from the database
// @Tags orders
// @Produce application/json
// @Success 200 {array} generated.Order
// @Router /api/orders [get]
func (h *OrderHandler) GetAllOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orders, err := h.domain.GetAllOrdersDomain(ctx)
		if err != nil {
			http.Error(w, "Failed to fetch restaurants", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(orders)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// GetOrderById godoc
//
// @Summary Get order by id
// @Description Fetches an order based on the id from the database
// @Tags orders
// @Produce application/json
// @Param id path string true "Order ID"
// @Success 200 {object} generated.Order
// @Router /api/orders/{id} [get]
func (h *OrderHandler) GetOrderById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orderIdStr := r.PathValue("orderId")
		if orderIdStr == "" {
			http.Error(w, "Missing orderId query parameter", http.StatusBadRequest)
			return
		}

		orderId, err := strconv.Atoi(orderIdStr)
		if err != nil {
			http.Error(w, "Invalid Order ID", http.StatusBadRequest)
			return
		}

		order, err := h.domain.GetOrderByIdDomain(ctx, int32(orderId))
		if err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(order)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" example:"Pending/On its way/Delivered"`
}

// UpdateOrderStatus godoc
//
// @Summary Update Order Status
// @Description Updates the status of an order
// @Tags orders
// @Accept application/json
// @Produce application/json
// @Param orderId path int true "Order ID"
// @Param status body UpdateOrderStatusRequest true "New Order Status"
// @Success 200 {string} string "Order status updated successfully"
// @Failure 404 {string} string "Order not found"
// @Router /api/order/status/{orderId} [patch]
func (h *OrderHandler) UpdateOrderStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orderIdStr := r.PathValue("orderId")
		orderId, err := strconv.Atoi(orderIdStr)
		if err != nil {
			http.Error(w, "Invalid Order ID", http.StatusBadRequest)
			return
		}

		// Parse the new status from the request body
		var requestPayload struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate the new status
		validStates := []string{"Pending", "On its way", "Delivered"}
		isValid := false
		for _, validStatus := range validStates {
			if requestPayload.Status == validStatus {
				isValid = true
				break
			}
		}
		if !isValid {
			http.Error(w, "Invalid status value, you can only choose between: Pending/On its way/Delivered", http.StatusBadRequest)
			return
		}

		// Call the domain fucntion to update the order status
		err = h.domain.UpdateOrderStatusDomain(ctx, int32(orderId), requestPayload.Status)
		if err != nil {
			if err.Error() == "order not found" {
				http.Error(w, "Order not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to update order status", http.StatusInternalServerError)
			}
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Order status updated successfully}`))
	}
}

type UpdateOrderStatusRequestWithDeliveryAgentId struct {
	DeliveryAgentId int32  `json:"id"`
	Status          string `json:"status" example:"Pending/On its way/Delivered"`
}

// UpdateOrderStatus godoc
//
// @Summary Update Order Status
// @Description Updates the status of an order
// @Tags orders
// @Accept application/json
// @Produce application/json
// @Param orderId path int true "Order ID"
// @Param status body UpdateOrderStatusRequestWithDeliveryAgentId true "New Order Status"
// @Success 200 {string} string "Order status updated successfully"
// @Failure 404 {string} string "Order not found"
// @Router /api/order/status-agent/{orderId} [patch]
func (h *OrderHandler) UpdateOrderStatusWithDeliveryAgentId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orderIdStr := r.PathValue("orderId")
		orderId, err := strconv.Atoi(orderIdStr)
		if err != nil {
			http.Error(w, "Invalid Order ID", http.StatusBadRequest)
			return
		}

		// Parse the new status from the request body
		var requestPayload struct {
			DeliveryAgentId int32  `json:"id"`
			Status          string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate the new status
		validStates := []string{"Pending", "On its way", "Delivered"}
		isValid := false
		for _, validStatus := range validStates {
			if requestPayload.Status == validStatus {
				isValid = true
				break
			}
		}
		if !isValid {
			http.Error(w, "Invalid status value, you can only choose between: Pending/On its way/Delivered", http.StatusBadRequest)
			return
		}

		// Call the domain fucntion to update the order status
		err = h.domain.UpdateOrderStatusAndDeliveryAgentDomain(ctx, int32(orderId), requestPayload.Status, requestPayload.DeliveryAgentId)
		if err != nil {
			if err.Error() == "order not found" {
				http.Error(w, "Order not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to update order status", http.StatusInternalServerError)
			}
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Order status updated successfully}`))
	}
}

// DeleteOrder godoc
//
// @Summary Delete an order
// @Description Deletes an order by its id from the database
// @Tags orders
// @Param id path int true "Order ID"
// @Success 200 {string} string "Order deleted successfully"
// @Failure 404 {string} string "Order not found"
// @Router /api/orders/{id} [delete]
func (h *OrderHandler) DeleteOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orderIdStr := r.PathValue("orderId")
		if orderIdStr == "" {
			http.Error(w, "Missing orderId query parameter", http.StatusBadRequest)
			return
		}

		orderId, err := strconv.Atoi(orderIdStr)
		if err != nil {
			http.Error(w, "Invalid Order ID", http.StatusBadRequest)
			return
		}

		err = h.domain.DeleteOrderDomain(ctx, int32(orderId))
		if err != nil {
			if err.Error() == "Order not found" {
				http.Error(w, "Order not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to delete order", http.StatusInternalServerError)
			}
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Order deleted successfully"}`))
	}
}

/* BROKER */

// Helper functions
func int32Ptr(i int) *int32 {
	value := int32(i)
	return &value
}
func toTimeNowPtr() *time.Time {
	now := time.Now()
	return &now
}

// ConsumeOrder godoc
//
//	@Summary		Consume Order for a Customer
//	@Description	Consume the created order for customer
//	@Tags			order
//	@Produce		application/json
//	@Success		200	{string}	string	"Order Consumed Successfully"
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/api/order/consume [get]
func (h *OrderHandler) ConsumeOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		broker.Consume("order_created_queue", func(event broker.Event) {
			// Ensure the event type is as expected
			if event.Type != broker.OrderCreated {
				log.Printf("Ignored event of unexpected type: %v", event.Type)
				return
			}

			// Convert event.Payload (interface{}) to JSON bytes
			payloadBytes, err := json.Marshal(event.Payload)
			if err != nil {
				log.Printf("Failed to marshal event payload: %v", err)
				return
			}

			// Unmarshal JSON into a struct making the Redis payload from ShoppingCart
			var payload struct {
				Customerid   int                               `json:"customer_id"`
				Restaurantid int                               `json:"restaurant_id"`
				Totalamount  float64                           `json:"total_amount"`
				Vatamount    float64                           `json:"vat_amount"`
				Comment      string                            `json:"comment"`
				Items        []generated.CreateOrderItemParams `json:"items"`
			}

			if err := json.Unmarshal(payloadBytes, &payload); err != nil {
				log.Printf("Failed to unmarshal payload into structured data: %v", err)
				return
			}
			log.Printf("Received payload: %+v", payload)

			// Create Order
			orderParams := generated.CreateOrderParams{
				Totalamount:     payload.Totalamount,
				Vatamount:       payload.Vatamount,
				Status:          "Pending",
				Timestamp:       toTimeNowPtr(),
				Comment:         &payload.Comment,
				Customerid:      int32Ptr(payload.Customerid),
				Restaurantid:    int32Ptr(payload.Restaurantid),
				Deliveryagentid: nil, // Not assigned yet
				Paymentid:       nil, // Not processed yet
				Bonusid:         nil, // No bonus assigned
				Feeid:           nil, // No fees applied
			}

			// Create context
			ctx := context.Background()

			// Call the CreateOrder domain function
			orderid, err := h.domain.CreateOrderDomain(ctx, orderParams)
			if err != nil {
				log.Printf("Failed to create order: %v", err)
				return
			}

			// Log success for the order creation
			log.Printf("Successfully created order with ID: %d for customer: %d", orderid, payload.Customerid)

			// Create order items for the created order
			for _, item := range payload.Items {
				itemParams := generated.CreateOrderItemParams{
					Orderid:  orderid,
					Name:     item.Name,
					Price:    item.Price,
					Quantity: item.Quantity,
				}

				// Call the CreateOrderItem domain function
				_, err := h.domain.CreateOrderItemDomain(ctx, itemParams)
				if err != nil {
					log.Printf("Failed to create order item: %+v, err: %v", item, err)
					continue
				}

				log.Printf("Successfully added item to order ID %d: %+v", orderid, item)
			}
		})

		// Respond to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Order consumption started",
		})
	}
}

// CalculateOrderBonus godoc
//
// @Summary calculate order bonus
// @Description calculates the order bonus
// @Tags orders
// @Param orderId path string true "Order ID"
// @Produce application/json
// @Success 200 {array} generated.Order
// @Router /api/order/bonus/{orderId} [get]
func (h *OrderHandler) CalculateOrderBonus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orderIdStr := r.PathValue("orderId")
		if orderIdStr == "" {
			http.Error(w, "Missing orderId query parameter", http.StatusBadRequest)
			return
		}

		orderId, err := strconv.Atoi(orderIdStr)
		if err != nil {
			http.Error(w, "Invalid Order ID", http.StatusBadRequest)
			return
		}

		totalBonus, err := h.domain.CalculateBonus(ctx, int32(orderId))
		if err != nil {
			http.Error(w, "Failed to calculate fee", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(totalBonus)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}
