package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rasm445f/soft-exam-2/broker"
	"github.com/rasm445f/soft-exam-2/domain"
)

type ShoppingCartHandler struct {
	domain domain.ShoppingCartPort
}

func NewShoppingCartHandler(domain domain.ShoppingCartPort) *ShoppingCartHandler {
	return &ShoppingCartHandler{domain: domain}
}

// AddItem godoc
//
//	@Summary		Add a MenuItem
//	@Description	Add a MenuItem to the shopping cart
//	@Tags			ShoppingCart CRUD
//	@Accept			application/json
//	@Produce		application/json
//	@Param			item	body		domain.AddItemParams	true	"item object"
//	@Success		201		{object}	domain.AddItemParams
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/api/shopping [post]
func (h *ShoppingCartHandler) AddItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var item domain.AddItemParams

		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := h.domain.AddItemDomain(ctx, item); err != nil {
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
//	@Tags			ShoppingCart CRUD
//	@Accept			application/json
//	@Produce		application/json
//	@Param			customerId	path		int						true	"customer ID"
//	@Param			itemId		path		int						true	"Item ID"
//	@Param			body		body		UpdateQuantityRequest	true	"New quantity for the item"
//	@Success		200			{string}	string					"ShoppingCart updated successfully"
//
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
//
//	@Router			/api/shopping/{customerId}/{itemId} [patch]
func (h *ShoppingCartHandler) UpdateCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		customerIdStr := r.PathValue("customerId")
		fmt.Printf("customerId: %v\n", customerIdStr)
		itemIdstr := r.PathValue("itemId")
		fmt.Printf("itemid: %v \n", itemIdstr)
		customerId, err := strconv.Atoi(customerIdStr)
		if err != nil {
			http.Error(w, "Malformed customer_id", http.StatusBadRequest)
		}
		itemId, err := strconv.Atoi(itemIdstr)
		if err != nil {
			http.Error(w, "Malformed item_id", http.StatusBadRequest)
		}

		var req UpdateQuantityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if err := h.domain.UpdateCartDomain(ctx, customerId, itemId, req.Quantity); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// ViewCart godoc
//
//	@Summary		View MenuItems for a customer's ShoppingCart
//	@Description	Fetches a list of MenuItems for a specific Customer, to view the ShoppingCart
//	@Tags			ShoppingCart CRUD
//	@Produce		application/json
//	@Param			id	path		string	true	"customer ID"
//	@Success		200	{string}	string	"Viewed Cart"
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

		shoppingCart, err := h.domain.ViewCartDomain(ctx, customerId)
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
//	@Summary		Clears the ShoppingCart
//	@Description	Clears the ShoppingCart for a specific customer
//	@Tags			ShoppingCart CRUD
//	@Accept			application/json
//	@Produce		application/json
//	@Param			customerId	path		int		true	"customer ID"
//	@Success		200			{string}	string	"Cart Cleared"
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

		if err := h.domain.ClearCartDomain(ctx, customerId); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// Consume Shopping Cart's MenuItems godoc
//
//	@Summary		Consume the chosen Menu Items for a Customer
//	@Description	Consumes the Shopping Cart's Menu Items for a Customer
//	@Tags			ShoppingCart Broker
//	@Produce		application/json
//	@Success		200	{string}	string	"Shopping Cart's Menu Items Consumed"
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

			// Unmarshal JSON bytes
			var item domain.AddItemParams
			if err := json.Unmarshal(payloadBytes, &item); err != nil {
				log.Printf("Failed to unmarshal payload: %v", err)
				return
			}
			fmt.Printf("Unmarshaled JSON bytes: %v", item)

			// Create a context for the AddItem function
			ctx := context.Background()

			// Call the AddItem logic
			if err := h.domain.AddItemDomain(ctx, item); err != nil {
				log.Printf("Failed to add MenuItem to shopping cart: %v", err)
				return
			}

			log.Printf("Successfully added MenuItem to shopping cart: %+v", item)
		})
	}
}

type PublishShoppingCartRequest struct {
	Comment string `json:"comment" example:"No vegetables on the pizza."`
}

// PublishShoppingCart godoc
//
//	@Summary		Publish a Customer's shopping cart to RabbitMQ to be consumed by the Order service with an optional Comment
//	@Description	Selecting the cart for the specified customer with an optional comment
//	@Tags			ShoppingCart Broker
//	@Accept			application/json
//	@Produce		application/json
//	@Param			customerId	path		int		true	"Customer ID"
//	@Param			comment		body		PublishShoppingCartRequest		true	"Customer Comment (optional)"
//	@Success		200			{string}	string	"Order Selected Successfully"
//	@Failure		400			{string}	string	"Bad request"
//	@Failure		500			{string}	string	"Internal server error"
//	@Router			/api/shopping/publish/{customerId} [post]
func (h *ShoppingCartHandler) PublishShoppingCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		customerIdStr := r.PathValue("customerId")
		customerId, err := strconv.Atoi(customerIdStr)
		if err != nil {
			http.Error(w, "Invalid customer_id", http.StatusBadRequest)
		}

		// Decode the Comment from the request body
		var requestPayload PublishShoppingCartRequest
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Fetch the shopping cart
		shoppingCart, err := h.domain.ViewCartDomain(ctx, customerId)
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
			http.Error(w, "Failed to publish shopping cart", http.StatusInternalServerError)
			return
		}
		// TODO: should the shoppingcart be cleared afterwards? or maybe when the order is confirmed?
		// h.domain.ClearCart(ctx, customerId)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Shopping Cart published and selected to Order successfully}`))
	}
}
