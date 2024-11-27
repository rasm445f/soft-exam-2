package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rasm445f/soft-exam-2/broker"
	"github.com/rasm445f/soft-exam-2/db"
	"github.com/rasm445f/soft-exam-2/domain"
)

type ShoppingCartHandler struct {
	domain *domain.ShoppingCartDomain
}

func NewShoppingCartHandler(domain *domain.ShoppingCartDomain) *ShoppingCartHandler {
	return &ShoppingCartHandler{domain: domain}
}

// AddItem godoc
//
//	@Summary		Add an item
//	@Description	Add an item to the shopping cart
//	@Tags			shoppingCart
//	@Accept			application/json
//	@Produce		application/json
//	@Param			item	body		db.AddItemParams	true	"item object"
//	@Success		201		{object}	db.AddItemParams
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/api/shopping [post]
func (h *ShoppingCartHandler) AddItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var item db.AddItemParams

		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := h.domain.AddItem(ctx, item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// Seperate struct needed for swagger documentation
type UpdateQuantityRequest struct {
	Quantity int `json:"quantity"`
}

// UpdateCartHandler godoc
//
//	@Summary		Update an item in the cart
//	@Description	Update the quantity of an existing item in the shopping cart. If the quantity is set to 0, the item will be removed.
//	@Tags			shoppingCart
//	@Accept			application/json
//	@Produce		application/json
//	@Param			customerId	path		int						true	"customer ID"
//	@Param			itemId		path		int						true	"Item ID"
//	@Param			body		body		UpdateQuantityRequest	true	"New quantity for the item"
//	@Success		200			{string}	string					"Item updated successfully"
//	@Failure		400			{string}	string					"Invalid input"
//	@Failure		404			{string}	string					"Item not found"
//	@Failure		500			{string}	string					"Internal server error"
//	@Router			/api/shopping/{customerId}/{itemId} [patch]
func (h *ShoppingCartHandler) UpdateCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		customerIdStr := r.PathValue("customerId")
		itemIdstr := r.PathValue("itemId")
		customerId, err := strconv.Atoi(customerIdStr)
		if err != nil {
			http.Error(w, "Malformed customer_id", http.StatusBadRequest)
		}
		itemId, err := strconv.Atoi(itemIdstr)
		if err != nil {
			http.Error(w, "Malformed customer_id", http.StatusBadRequest)
		}

		var req UpdateQuantityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if err := h.domain.UpdateCart(ctx, customerId, itemId, req.Quantity); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// ViewCart godoc
//
//	@Summary		View items for a customer
//	@Description	Fetches a list of items based on the customerId
//	@Tags			shoppingCart
//	@Produce		application/json
//	@Param			id	path		string	true	"customer ID"
//	@Success		200	{string}	string	"Cart cleared"
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/api/shopping/{id} [get]
func (h *ShoppingCartHandler) ViewCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		customerIdStr := r.PathValue("customerId")
		customerId, err := strconv.Atoi(customerIdStr)
		if err != nil {
			http.Error(w, "Malformed customer_id", http.StatusBadRequest)
		}

		shoppingCart, err := h.domain.ViewCart(ctx, customerId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(shoppingCart); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// ClearCart godoc
//
//	@Summary		Clears the cart
//	@Description	Clears the cart for the specified customer
//	@Tags			shoppingCart
//	@Accept			application/json
//	@Produce		application/json
//	@Param			customerId	path		int		true	"customer ID"
//	@Success		200			{string}	string	"cart cleared"
//	@Failure		400			{string}	string	"Bad request"
//	@Failure		500			{string}	string	"Internal server error"
//	@Router			/api/shopping/{customerId} [delete]
func (h *ShoppingCartHandler) ClearCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		customerIdStr := r.PathValue("customerId")
		customerId, err := strconv.Atoi(customerIdStr)
		if err != nil {
			http.Error(w, "Malformed customer_id", http.StatusBadRequest)
		}

		if err := h.domain.ClearCart(ctx, customerId); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// Consume godoc
//
//	@Summary		View items for a customer
//	@Description	Fetches a list of items based on the customerId
//	@Tags			shoppingCart
//	@Produce		application/json
//	@Success		200	{string}	string	"Cart cleared"
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/api/shopping/consume [get]
func (h *ShoppingCartHandler) ConsumeMenuItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		broker.Consume("menu_item_selected_queue", func(event broker.Event) {
			// Ensure the event type is as expected
			if event.Type != broker.MenuItemSelected {
				log.Printf("Ignored event of unexpected type: %v", event.Type)
				return
			}

			// Convert event.Payload (interface{}) to JSON bytes
			payloadBytes, err := json.Marshal(event.Payload)
			if err != nil {
				log.Printf("Failed to marshal event payload: %v", err)
				return
			}
			fmt.Println(payloadBytes)

			// Unmarshal JSON bytes into IntermediatePayload
			var item db.AddItemParams
			if err := json.Unmarshal(payloadBytes, &item); err != nil {
				log.Printf("Failed to unmarshal payload into IntermediatePayload: %v", err)
				return
			}
			fmt.Printf("Unmarshaled intermediate: %v", item)

			// Create a context for the AddItem function
			ctx := context.Background()

			// Call the AddItem logic
			if err := h.domain.AddItem(ctx, item); err != nil {
				log.Printf("Failed to add item to shopping cart: %v", err)
				return
			}

			log.Printf("Successfully added item to shopping cart: %+v", item)
		})
	}
}

// ClearCart godoc
//
//	@Summary		Publish the shopping cart to rabbimq to be consumed by the Order service
//	@Description	Clears the cart for the specified customer
//	@Tags			shoppingCart
//	@Accept			application/json
//	@Produce		application/json
//	@Param			customerId	path		int		true	"customer ID"
//	@Success		200			{string}	string	"cart cleared"
//	@Failure		400			{string}	string	"Bad request"
//	@Failure		500			{string}	string	"Internal server error"
//	@Router			/api/shopping/publish/{customerId} [get]
func (h *ShoppingCartHandler) SelectOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		customerIdStr := r.PathValue("customerId")
		customerId, err := strconv.Atoi(customerIdStr)
		if err != nil {
			http.Error(w, "Malformed customer_id", http.StatusBadRequest)
		}
		shoppingCart, err := h.domain.ViewCart(ctx, customerId)
		if err != nil {
			log.Printf("Failed to publish event: %v", err)
			http.Error(w, "Failed to select menu item", http.StatusInternalServerError)
			return
		}

		// Publish event to RabbitMQ
		event := broker.Event{
			Type:    broker.OrderCreated,
			Payload: shoppingCart,
		}
		err = broker.Publish("order_created_queue", event)
		if err != nil {
			log.Printf("Failed to publish event: %v", err)
			http.Error(w, "Failed to select menu item", http.StatusInternalServerError)
			return
		}
		// TODO: should the shoppingcart be cleared afterwards? or maybe when the order is confirmed?
		// h.domain.ClearCart(ctx, customerId)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Menu item selected successfully}`))
	}
}
