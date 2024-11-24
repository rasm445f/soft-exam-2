package handlers

import (
	"log"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/oTuff/go-startkode/db"
<<<<<<< HEAD
	"github.com/rasm445f/soft-exam-2/broker"
)

type ShoppingCart struct {
	CustomerID string `json:"customerId"`
	Items []db.ShoppingCartItem `json:"items"`
}

// PublishUpdatedCart publishes the updated shopping cart to RabbitMQ
func PublishUpdatedCart(cart ShoppingCart) {
	event := broker.Event{
		Type: broker.CartUpdated,
		Payload: cart,
	}
	err := broker.Publish("cart_updated_queue", event)
		if err != nil {
			log.Printf("Failed to publish cart update: %v", err)
		}
}

// AddItemHandler processes menu_item_selected events and updates the cart
// AddItemHandler godoc
// @Summary Add an item
// @Description Add an item to the shopping cart
// @Tags shoppingCart
// @Accept  application/json
// @Produce application/json
// @Param item body db.AddItemParams true "item object"
// @Success 201 {object} db.AddItemParams
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/shopping [post]
func AddItemHandler(commands *db.ShoppingCartRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var selection broker.Event

		// Decode the RabbitMQ eventq
		err := json.NewDecoder(r.Body).Decode(&selection)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add item to the shopping cart
		var cartItem db.AddItemParams
		json.Unmarshal(selection.Payload, &cartItem)

		err = commands.AddItem(ctx, cartItem)
		if err != nil {
=======
	"github.com/oTuff/go-startkode/domain"
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
>>>>>>> origin/main
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Fetch the updated cart
		cart, err := commands.ViewCart(ctx, cartItem.CustomerID)
		if err != nil {
			http.Error(w, "Failed to fetch updated cart", http.StatusInternalServerError)
			return
		}

		// Publish the updated cart to RabbitMQ
		updatedCart := ShoppingCart{
			CustomerID: cartItem.CustomerID,
			Items: cart,
		}
		PublishUpdatedCart(updatedCart)

		w.WriteHeader(http.StatusOK)
	}
}

<<<<<<< HEAD

// AddItemHandler godoc
// @Summary Add an item
// @Description Add an item to the shopping cart
// @Tags shoppingCart
// @Accept  application/json
// @Produce application/json
// @Param item body db.AddItemParams true "item object"
// @Success 201 {object} db.AddItemParams
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/shopping [post]
// func AddItemHandler(commands *db.ShoppingCartRepository) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		var shoppingCartItem db.AddItemParams

// 		err := json.NewDecoder(r.Body).Decode(&shoppingCartItem)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		err = commands.AddItem(ctx, shoppingCartItem)

// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		w.WriteHeader(http.StatusCreated)
// 	}
// }

// GetTodo godoc
// @Summary View items for a user
// @Description Fetches a list of items based on the userId
// @Tags shoppingCart
// @Produce application/json
// @Param id path string true "User ID"
// @Success 200 {array} db.ShoppingCartItem
// @Router /api/shopping/{id} [get]
func ViewItemHandler(commands *db.ShoppingCartRepository) http.HandlerFunc {
=======
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
>>>>>>> origin/main
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
