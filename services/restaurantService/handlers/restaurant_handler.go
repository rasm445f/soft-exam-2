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

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		restaurant, err := queries.GetRestaurantById(ctx, int32(id))
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

func GetMenuItemsByRestaurant(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		restaurantIdStr := r.URL.Query().Get("restaurantId")
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