package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/oTuff/go-startkode/db/generated"
	"github.com/oTuff/go-startkode/domain"
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

		err = h.domain.CreateCustomer(ctx, customer)
		if err != nil {
			http.Error(w, "Failed to create customer", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

// UpdateCustomer godoc
// @Summary Update a customer
// @Description Updates a customer's details based on the ID provided in the URL path. This may include personal information as well as optional address updates.
// @Tags customers
// @Accept application/json
// @Produce application/json
// @Param id path int true "Customer ID"
// @Param customer body generated.UpdateCustomerParams true "Updated customer details"
// @Success 200 {object} map[string]string "Customer updated successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Customer not found"
// @Failure 500 {object} map[string]string "Internal server error"
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

		// Decode the incoming JSON request into a map to capture all fields
		var updatePayload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updatePayload); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Println("Error decoding request body:", err)
			return
		}

		// Create an UpdateCustomerParams struct and fill it based on the JSON payload
		customerUpdates := generated.UpdateCustomerParams{
			ID: int32(id),
		}
		if name, ok := updatePayload["name"].(string); ok {
			customerUpdates.Name = &name
		}
		if email, ok := updatePayload["email"].(string); ok {
			customerUpdates.Email = &email
		}
		if phoneNumber, ok := updatePayload["phonenumber"].(string); ok {
			customerUpdates.Phonenumber = &phoneNumber
		}
		if password, ok := updatePayload["password"].(string); ok {
			customerUpdates.Password = &password
		}

		// Update the customer information in the database
		err = h.domain.UpdateCustomer(ctx, customerUpdates)
		if err != nil {
			if err.Error() == "customer not found" {
				http.Error(w, "Customer not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to update customer", http.StatusInternalServerError)
			}
			log.Println("Error updating customer:", err)
			return
		}

		// Check if the address needs to be updated
		if streetAddress, ok := updatePayload["street_address"].(string); ok || updatePayload["zip_code"] != nil {
			addressUpdates := generated.UpdateAddressParams{
				ID: int32(id),
			}

			if ok {
				addressUpdates.StreetAddress = &streetAddress
			}
			if zipCode, ok := updatePayload["zip_code"].(int32); ok {
				addressUpdates.ZipCode = &zipCode
			}

			// Update the address in the database
			err = h.domain.UpdateAddress(ctx, addressUpdates)
			if err != nil {
				http.Error(w, "Failed to update address", http.StatusInternalServerError)
				log.Println("Error updating address:", err)
				return
			}
		}

		// Return a success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Customer updated successfully"}`))
	}
}
