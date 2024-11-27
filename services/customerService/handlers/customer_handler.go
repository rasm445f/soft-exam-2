package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/rasm445f/soft-exam-2/broker"
	"github.com/rasm445f/soft-exam-2/db/generated"
	"github.com/rasm445f/soft-exam-2/domain"
)

type CustomerHandler struct {
	domain *domain.CustomerDomain
}

func NewCustomerHandler(domain *domain.CustomerDomain) *CustomerHandler {
	return &CustomerHandler{domain: domain}
}

// GetAllCustomers godoc
// @Summary Get all customers
// @Description Fetches a list of all customers from the database
// @Tags customers
// @Produce application/json
// @Success 200 {array} generated.Customer
// @Router /api/customer [get]
func (h *CustomerHandler) GetAllCustomers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		customers, err := h.domain.GetAllCustomers(ctx)
		if err != nil {
			http.Error(w, "Failed to fetch customers", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(customers)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

// GetCustomerById godoc
// @Summary Get customer
// @Description Fetches a customer based on the id from the database
// @Tags customers
// @Produce application/json
// @Param id path string true "Customer ID"
// @Success 200 {object} generated.Customer
// @Router /api/customer/{id} [get]
func (h *CustomerHandler) GetCustomerById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			log.Println(err)
			return
		}

		customer, err := h.domain.GetCustomerByID(ctx, int32(id))
		if err != nil {
			http.Error(w, "Customer not found", http.StatusNotFound)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(customer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)

	}
}

// DeleteCustomer godoc
// @Summary Delete customer
// @Description Deletes a customer based on the id from the database
// @Tags customers
// @Produce application/json
// @Param id path string true "Customer ID"
// @Success 200 {string} string "Customer deleted"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Customer not found"
// @Router /api/customer/{id} [delete]
func (h *CustomerHandler) DeleteCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			log.Println(err)
			return
		}

		err = h.domain.DeleteCustomer(ctx, int32(id))
		if err != nil {
			http.Error(w, "Customer not found", http.StatusNotFound)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Customer deleted"}`))
	}
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Creates a new customer entry in the database
// @Tags customers
// @Accept  application/json
// @Produce application/json
// @Param customer body generated.CreateCustomerParams true "Customer object"
// @Success 201 {object} generated.Customer
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/customer [post]
func (h *CustomerHandler) CreateCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var customer generated.CreateCustomerParams

		err := json.NewDecoder(r.Body).Decode(&customer)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Println(err)
			return
		}

		if *customer.Name == "" || *customer.Email == "" || *customer.Password == "" {
			http.Error(w, "All required fields must be filled", http.StatusBadRequest)
			return
		}

		createdCustomer, err := h.domain.CreateCustomer(ctx, customer)
		if err != nil {
			http.Error(w, "Failed to create customer", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res, _ := json.Marshal(createdCustomer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}
}

// UpdateCustomer godoc
// @Summary Update customer
// @Description Updates a customer's details based on the ID from the database
// @Tags customers
// @Accept application/json
// @Produce application/json
// @Param id path string true "Customer ID"
// @Param customer body generated.UpdateCustomerParams true "Updated customer details"
// @Success 200 {string} string "Customer updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Customer not found"
// @Router /api/customer/{id} [patch]
func (h *CustomerHandler) UpdateCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			log.Println("Error parsing customer ID:", err)
			return
		}

		// Decode the incoming JSON request
		var customerUpdates generated.UpdateCustomerParams
		if err := json.NewDecoder(r.Body).Decode(&customerUpdates); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Println("Error decoding request body:", err)
			return
		}

		// Ensure the ID matches the customer's ID being updated
		customerUpdates.ID = int32(id)

		// Call the query to update the customer in the database
		err = h.domain.UpdateCustomer(ctx, customerUpdates)
		if err != nil {
			if err.Error() == "Customer not found" {
				http.Error(w, "Customer not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to update customer", http.StatusInternalServerError)
			}
			log.Println("Error updating customer:", err)
			return
		}

		// Return a success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Customer updated successfully"}`))
	}
}

/* BROKER */
type MenuItemSelection struct {
	CustomerID   int32 `json:"customerId" example:"1"`
	RestaurantId int32 `json:"restaurantId" example:"10"`
	Name   string `json:"name" example:"Cheese Burger"`
	Price  float64 `json:"price" example:"9.99"`
	Quantity     int   `json:"quantity" example:"2"`
}
// SelectMenuItem godoc
// @Summary Select Menuitem
// @Description Select Menu Item
// @Tags customers
// @Accept  application/json
// @Produce application/json
// @Param customer body MenuItemSelection true "Menu item selection details"
// @Success 201 {object} MenuItemSelection "Menu item successfully selected"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/customer/menu/select [post]
func (h *CustomerHandler) SelectMenuitem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var selection MenuItemSelection
		err := json.NewDecoder(r.Body).Decode(&selection)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Publish event to RabbitMQ
		event := broker.Event{
			Type:    broker.MenuItemSelected,
			Payload: selection,
		}
		err = broker.Publish("menu_item_selected_queue", event)
		if err != nil {
			log.Printf("Failed to publish event: %v", err)
			http.Error(w, "Failed to select menu item", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Menu item selected successfully}`))
	}
}
