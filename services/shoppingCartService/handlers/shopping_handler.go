package handlers

import (
	"log"
	"encoding/json"
	"net/http"

	"github.com/oTuff/go-startkode/db"
	"github.com/oTuff/go-startkode/broker"
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
func AddItemHandler(commands *db.ShoppingCartRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var selection broker.Event

		// Decode the RabbitMQ event
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
		var shoppingCartItem db.AddItemParams

		err := json.NewDecoder(r.Body).Decode(&shoppingCartItem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = commands.AddItem(ctx, shoppingCartItem)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// GetTodo godoc
// @Summary View items for a user
// @Description Fetches a list of items based on the userId
// @Tags shoppingCart
// @Produce application/json
// @Param id path string true "User ID"
// @Success 200 {array} db.ShoppingCartItem
// @Router /api/shopping/{id} [get]
func ViewItemHandler(commands *db.ShoppingCartRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.PathValue("userId")
		items, err := commands.ViewCart(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, _ := json.Marshal(items)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)

	}
}
