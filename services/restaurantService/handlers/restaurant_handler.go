package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/rasm445f/soft-exam-2/db/generated"
)

// GetAllRestaurants godoc
// @Summary Get all restaurants
// @Description Fetches a list of all restaurants from the database
// @Tags restaurants
// @Produce application/json
// @Success 200 {array} generated.Restaurant
// @Router /api/restaurant [get]
func GetAllRestaurants(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		restaurants, err := queries.FetchAllRestaurants(ctx)
		if err != nil {
			http.Error(w, "Failed to fetch restaurants", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(restaurants)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// GetRestaurantById godoc
// @Summary Get restaurant by id
// @Description Fetches a restaurant based on the id from the database
// @Tags restaurants
// @Produce application/json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} generated.Restaurant
// @Router /api/restaurant/{id} [get]
func GetRestaurantById(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		restaurantIdStr := r.PathValue("restaurantId")
		if restaurantIdStr == "" {
			http.Error(w, "Missing restaurantId query parameter", http.StatusBadRequest)
			return
		}

		restaurantId, err := strconv.Atoi(restaurantIdStr)
		if err != nil {
			http.Error(w, "Invalid Restaurant ID", http.StatusBadRequest)
			return
		}

		restaurant, err := queries.GetRestaurantById(ctx, int32(restaurantId))
		if err != nil {
			http.Error(w, "Restaurant not found", http.StatusNotFound)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(restaurant)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// GetMenuItemsByRestaurant godoc
// @Summary Get menu items by restaurant ID
// @Description Fetches all menu items associated with a specific restaurant ID
// @Tags menu_items
// @Produce application/json
// @Param restaurantId query string true "Restaurant ID"
// @Success 200 {array} generated.MenuItem
// @Failure 400 {string} string "Invalid Rastaurant ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/restaurant/{restaurantId}/menu-items [get]
func GetMenuItemsByRestaurant(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// restaurantIdStr := r.URL.Path[len("/api/restaurants/") : len(r.URL.Path)-len("/menu-items")]
		restaurantIdStr := r.PathValue("restaurantId")
		if restaurantIdStr == "" {
			http.Error(w, "Missing restaurantId path parameter", http.StatusBadRequest)
			return
		}

		restaurantId, err := strconv.Atoi(restaurantIdStr)
		if err != nil {
			http.Error(w, "Invalid Restaurant ID", http.StatusBadRequest)
			return
		}

		menuItems, err := queries.FetchMenuItemsByRestaurantId(ctx, int32(restaurantId))
		if err != nil {
			http.Error(w, "Failed to fetch menu items", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(menuItems)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// GetMenuItemByRestaurantAndId godoc
// @Summary Get menu item by restaurant and id
// @Description Fetches a menu item based on the restaurant and id from the database
// @Tags menu_items
// @Produce application/json
// @Param id path string true "Menu Item ID"
// @Success 200 {object} generated.MenuItem
// @Router /api/restaurant/{restaurantId}/menu-items/{menuitemId} [get]
func GetMenuItemByRestaurantAndId(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		restaurantIdStr := r.PathValue("restaurantId")
		menuitemIdStr := r.PathValue("menuitemId")

		if restaurantIdStr == "" || menuitemIdStr == "" {
			http.Error(w, "Missing path parameters (restaurantId, menuitemId)", http.StatusBadRequest)
			return
		}

		restaurantId, err := strconv.Atoi(restaurantIdStr)
		if err != nil {
			http.Error(w, "Invalid Restaurant ID", http.StatusBadRequest)
			return
		}
		menuitemId, err := strconv.Atoi(menuitemIdStr)
		if err != nil {
			http.Error(w, "Invalid Restaurant ID", http.StatusBadRequest)
			return
		}

		params := generated.GetMenuItemByRestaurantAndIdParams{
			Restaurantid: int32(restaurantId),
			ID:           int32(menuitemId),
		}

		menuItem, err := queries.GetMenuItemByRestaurantAndId(ctx, params)
		if err != nil {
			http.Error(w, "Menu Item not found", http.StatusNotFound)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(menuItem)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}
