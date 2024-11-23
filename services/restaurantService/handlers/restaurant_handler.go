package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/rasm445f/soft-exam-2/db/generated"
	"github.com/rasm445f/soft-exam-2/domain"
)

type RestaurantHandler struct {
	domain *domain.RestaurantDomain
}

func NewRestaurantHandler(domain *domain.RestaurantDomain) *RestaurantHandler {
	return &RestaurantHandler{domain: domain}
}

// GetAllRestaurants godoc
// @Summary Get all restaurants
// @Description Fetches a list of all restaurants from the database
// @Tags restaurants
// @Produce application/json
// @Success 200 {array} generated.Restaurant
// @Router /api/restaurants [get]
func (h *RestaurantHandler) GetAllRestaurants() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		restaurants, err := h.domain.FetchAllRestaurants(ctx)
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
// @Router /api/restaurants/{id} [get]
func (h *RestaurantHandler) GetRestaurantById() http.HandlerFunc {
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

		restaurant, err := h.domain.GetRestaurantById(ctx, int32(restaurantId))
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
// @Param restaurantId path string true "Restaurant ID"
// @Success 200 {array} generated.Menuitem
// @Router /api/restaurants/{restaurantId}/menu-items [get]
func (h *RestaurantHandler) GetMenuItemsByRestaurant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

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

		menuItems, err := h.domain.FetchMenuItemsByRestaurantId(ctx, int32(restaurantId))
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
// @Param restaurantId path string true "Restaurant ID"
// @Param menuitemId path string true "Menu Item ID"
// @Success 200 {object} generated.Menuitem
// @Router /api/restaurants/{restaurantId}/menu-items/{menuitemId} [get]
func (h *RestaurantHandler) GetMenuItemByRestaurantAndId() http.HandlerFunc {
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

		menuItem, err := h.domain.GetMenuItemByRestaurantAndId(ctx, params)
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


// Categories

// GetAllCategories godoc
// @Summary Get all categories
// @Description Fetches a list of all unique categories from the restaurant
// @Tags categories
// @Produce application/json
// @Success 200 {array} string
// @Router /api/categories [get]
func (h *RestaurantHandler) GetAllCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		categories, err := h.domain.FetchAllCategories(ctx)
		if err != nil {
			http.Error(w, "Failed to fetch restaurants", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(categories)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// FilterRestaurantByCategory godoc
// @Summary Filter restaurants by category
// @Description Fetches all restaurants for a given category
// @Tags categories
// @Produce application/json
// @Param category path string true "Restaurant Category"
// @Success 200 {array} generated.Restaurant
// @Router /api/filter/{category} [get]
func (h *RestaurantHandler) FilterRestaurantByCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		category := r.PathValue("category")
		if category == "" {
			http.Error(w, "Missing category path parameter", http.StatusBadRequest)
			return
		}

		restaurants, err := h.domain.FilterRestaurantsByCategory(ctx, category)
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