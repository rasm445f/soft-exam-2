package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/oTuff/go-startkode/db"
)

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
