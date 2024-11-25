package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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


type IntermediatePayload struct {
	CustomerId   int `json:"customerId"`
	MenuItemId   int `json:"menuItemId"`
	Quantity     int `json:"quantity"`
	RestaurantId int `json:"restaurantId"`
}

// Consume godoc
//
//	@Summary		View order for a customer
//	@Description	Fetches a list of items based on the order
//	@Tags			order
//	@Produce		application/json
//	@Success		200	{string}	string	"Order Consume Success"
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
				Customerid int `json:"customer_id"`
				Restaurantid int `json:"restaurant_id"`
				Totalamount float64 `json:"total_amount"`
				Vatamount float64 `json:"vat_amount"`
				Comment string `json:"comment"`
				Items []generated.CreateOrderItemParams `json:"items"`
			}

			if err := json.Unmarshal(payloadBytes, &payload); err != nil {
				log.Printf("Failed to unmarshal payload into structured data: %v", err)
				return
			}
			log.Printf("Received payload: %+v", payload)

			// Create Order
			orderParams := generated.CreateOrderParams{
				Totalamount: payload.Totalamount,
				Vatamount: payload.Vatamount,
				Status: "Pending",
				Comment: &payload.Comment,
				
			}

			// Create a context for the AddItem function
			ctx := context.Background()

			// Call the AddItem logic
			if err := h.domain.CreateOrder(ctx, ); err != nil {
				log.Printf("Failed to add item to shopping cart: %v", err)
				return
			}

			log.Printf("Successfully added item to shopping cart: %+v", item)
		})
	}
}
