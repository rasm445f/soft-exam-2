package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/rasm445f/soft-exam-2/db/generated"
	"github.com/rasm445f/soft-exam-2/domain"
)

type CustomerHandler struct {
	domain domain.CustomerPort
}

func NewCustomerHandler(domain domain.CustomerPort) *CustomerHandler {
	return &CustomerHandler{domain: domain}
}

// GetAllCustomers godoc
//
// @Summary Get all customers
// @Description Fetches a list of all customers from the database
// @Tags Customer CRUD
// @Produce application/json
// @Success 200 {array} generated.Customer
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/customer [get]
func (h *CustomerHandler) GetAllCustomers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		customers, err := h.domain.GetAllCustomersDomain(ctx)
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
//
// @Summary Get customer by Id
// @Description Fetches a customer based on the id from the database
// @Tags Customer CRUD
// @Produce application/json
// @Param id path string true "Customer ID"
// @Success 200 {object} generated.Customer
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
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

		customer, err := h.domain.GetCustomerByIdDomain(ctx, int32(id))
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
//
// @Summary Delete customer
// @Description Deletes a customer based on the id from the database
// @Tags Customer CRUD
// @Produce application/json
// @Param id path string true "Customer ID"
// @Success 200 {string} string "Customer deleted"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
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

		err = h.domain.DeleteCustomerDomain(ctx, int32(id))
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
//
// @Summary Create a new customer
// @Description Creates a new customer entry in the database
// @Tags Customer CRUD
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

		err = h.domain.CreateCustomerDomain(ctx, customer)
		if err != nil {
			http.Error(w, "Failed to create customer", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

// UpdateCustomerWithAddress represents a customer update with optional address fields.
// @Description Update customer details including name, email, and address information.
type UpdateCustomerWithAddress struct {
	Name          *string `json:"name" example:"John Doe"`
	Email         *string `json:"email" example:"john.doe@example.com"`
	Phonenumber   *string `json:"phonenumber" example:"1234567890"`
	Password      *string `json:"password" example:"Password123!"`
	StreetAddress *string `json:"street_address" example:"123 Main St" `
	ZipCode       *int32  `json:"zip_code" example:"12345"`
	City          *string `json:"city" example:"New York"`
}

// UpdateCustomer godoc
//
// @Summary Update a customer
// @Description Updates a customer's details based on the ID provided in the URL path. This may include personal information as well as optional address updates.
// @Tags Customer CRUD
// @Accept application/json
// @Produce application/json
// @Param id path int true "Customer ID"
// @Param customer body UpdateCustomerWithAddress true "Updated customer details"
// @Success 200 {object} map[string]string "Customer updated successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/customer/{id} [patch]
func (h *CustomerHandler) UpdateCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr := r.PathValue("id")                 //source of tainted data
		id, err := strconv.ParseInt(idStr, 10, 64) //works as validation of tainted data
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			log.Println("Error parsing customer ID:", err)
			return
		}

		// Decode the incoming JSON request into a map to capture all fields
		var updatePayload map[string]interface{}                               //source of tainted data
		if err := json.NewDecoder(r.Body).Decode(&updatePayload); err != nil { //works as validation of tainted data
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

		// Call the query to update the customer in the database
		err = h.domain.UpdateCustomerDomain(ctx, customerUpdates)
		// Update the customer information in the database
		if err != nil {
			if err == sql.ErrNoRows {
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
			if zipCode, ok := updatePayload["zip_code"].(float64); ok {
				zipCodeInt32 := int32(zipCode)
				addressUpdates.ZipCode = &zipCodeInt32
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
		w.Write([]byte(`{"message": "Customer updated successfully"}`)) //returning json data therfore not a sink
	}
}
