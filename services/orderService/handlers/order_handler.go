package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
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

// Helper function
func int32Ptr(i int) *int32 {
	value := int32(i)
	return &value
}

// Consume godoc
//
//	@Summary		View order for a customer
//	@Description	Consume the created order for customer
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
				Totalamount pgtype.Numeric `json:"total_amount"`
				Vatamount pgtype.Numeric `json:"vat_amount"`
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
				Customerid: int32Ptr(payload.Customerid),
				Restaurantid: int32Ptr(payload.Restaurantid),
				Deliveryagentid: nil,	// Not assigned yet
				Paymentid: nil,			// Not processed yet
				Bonusid: nil,			// No bonus assigned
				Feeid: nil,				// No fees applied
			}

			// Create context
			ctx := context.Background()

			// Call the CreateOrder domain function
			orderid, err := h.domain.CreateOrder(ctx, orderParams)
			if err != nil {
				log.Printf("Failed to create order: %v", err)
				return
			}

			// Log success for the order creation
			log.Printf("Successfully created order with ID: %d for customer: %d", orderid, payload.Customerid)

			// Create order items for the created order
			for _, item := range payload.Items {
				itemParams := generated.CreateOrderItemParams{
					Orderid: orderid,
					Name: item.Name,
					Price: item.Price,
					Quantity: item.Quantity,
				}
				
				// Call the CreateOrderItem domain function
				_, err := h.domain.CreateOrderItem(ctx, itemParams)
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
