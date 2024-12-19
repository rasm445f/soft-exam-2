package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	brokerrestaurant "github.com/rasm445f/soft-exam-2/brokerRestaurant"
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
//
// @Summary Get all restaurants
// @Description Fetches a list of all restaurants from the database
// @Tags Restaurant CRUD
// @Produce application/json
// @Success 200 {array} generated.Restaurant
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/restaurants [get]
func (h *RestaurantHandler) GetAllRestaurants() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		restaurants, err := h.domain.GetAllRestaurantsDomain(ctx)
		if err != nil {
			http.Error(w, "Failed to get restaurants", http.StatusInternalServerError)
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
//
// @Summary Get restaurant by id
// @Description Fetches a restaurant based on the id from the database
// @Tags Restaurant CRUD
// @Produce application/json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} generated.Restaurant
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
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

		restaurant, err := h.domain.GetRestaurantByIdDomain(ctx, int32(restaurantId))
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
//
// @Summary Get menu items by restaurant ID
// @Description Fetches all menu items associated with a specific restaurant ID
// @Tags MenuItem(Restaurant) CRUD
// @Produce application/json
// @Param restaurantId path string true "Restaurant ID"
// @Success 200 {array} generated.Menuitem
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
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

		menuItems, err := h.domain.GetMenuItemsByRestaurantIdDomain(ctx, int32(restaurantId))
		if err != nil {
			http.Error(w, "Failed to get menu items", http.StatusInternalServerError)
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
//
// @Summary Get menu item by restaurant and id
// @Description Fetches a menu item based on the restaurant and id from the database
// @Tags MenuItem(Restaurant) CRUD
// @Produce application/json
// @Param restaurantId path string true "Restaurant ID"
// @Param menuitemId path string true "Menu Item ID"
// @Success 200 {object} generated.Menuitem
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
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

		menuItem, err := h.domain.GetMenuItemByRestaurantAndIdDomain(ctx, params)
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

/* CATEGORIES */

// GetAllCategories godoc
//
// @Summary Get all categories
// @Description Fetches a list of all unique categories from the restaurant
// @Tags Category(Restaurant) CRUD
// @Produce application/json
// @Success 200 {array} string
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/categories [get]
func (h *RestaurantHandler) GetAllCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		categories, err := h.domain.GetAllCategoriesDomain(ctx)
		if err != nil {
			http.Error(w, "Get to fetch restaurants", http.StatusInternalServerError)
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
//
// @Summary Filter restaurants by category
// @Description Fetches all restaurants for a given category
// @Tags Category(Restaurant) CRUD
// @Produce application/json
// @Param category path string true "Restaurant Category"
// @Success 200 {array} generated.Restaurant
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/filter/{category} [get]
func (h *RestaurantHandler) FilterRestaurantByCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		category := r.PathValue("category")
		if category == "" {
			http.Error(w, "Missing category path parameter", http.StatusBadRequest)
			return
		}

		restaurants, err := h.domain.FilterRestaurantsByCategoryDomain(ctx, category)
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

// SelectMenuItem godoc
//
// @Summary Selecting MenuItems
// @Description Customer selects a MenuItem or more
// @Tags Restaurant Broker
// @Accept  application/json
// @Produce application/json
// @Param customer body SelectItemParams true "Menu item selection details"
// @Success 201 {object} MenuItemSelection "Menu item successfully selected"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/restaurants/menu/select [post]
func (h *RestaurantHandler) SelectMenuItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var selectionParams brokerrestaurant.SelectItemParams
		err := json.NewDecoder(r.Body).Decode(&selectionParams)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		err = brokerrestaurant.SelectMenuBroker(selectionParams, ctx, h.domain)
		// if err != nil {
		// 	http.Error(w, "Menu Item not found", http.StatusNotFound)
		// 	log.Println(err)
		// 	return
		// }
		if err != nil {
			log.Printf("Failed to publish event: %v", err)
			http.Error(w, "Failed to select menu item", http.StatusInternalServerError)
			return
		}
		// Get menuItem based on restaurantId and menuItemId

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Menu item selected successfully}`))
	}
}
