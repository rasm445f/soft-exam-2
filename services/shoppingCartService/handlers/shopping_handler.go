package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/oTuff/go-startkode/db"
	"github.com/oTuff/go-startkode/domain"
)

type ShoppingCartHandler struct {
	domain *domain.ShoppingCartDomain
}

func NewShoppingCartHandler(domain *domain.ShoppingCartDomain) *ShoppingCartHandler {
	return &ShoppingCartHandler{domain: domain}
}

// AddItemHandler godoc
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
//	@Param			itemId	path		int						true	"Item ID"
//	@Param			body	body		UpdateQuantityRequest	true	"New quantity for the item"
//	@Success		200		{string}	string					"Item updated successfully"
//	@Failure		400		{string}	string					"Invalid input"
//	@Failure		404		{string}	string					"Item not found"
//	@Failure		500		{string}	string					"Internal server error"
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

// GetTodo godoc
//
//	@Summary		View items for a customer
//	@Description	Fetches a list of items based on the customerId
//	@Tags			shoppingCart
//	@Produce		application/json
//	@Param			id	path	string	true	"customer ID"
//	@Success		200	{array}	db.ShoppingCartItem
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

// func (h *ShoppingCartHandler) ClearCart() http.HandlerFunc {
// }
