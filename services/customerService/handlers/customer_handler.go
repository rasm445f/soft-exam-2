package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/oTuff/go-startkode/db/generated"
	"github.com/oTuff/go-startkode/mailer"
	"github.com/rasm445f/soft-exam-2/broker"
)

type MenuItemSelection struct {
	CustomerID int32 `json:"customerId"`
	MenuItemId int32 `json:"menuItemId"`
	Quantity int `json:"quantity"`
}


// SelectMenuItem godoc
// @Summary Select Menuitem
// @Description Select Menu Item
// @Tags customers
// @Accept  application/json
// @Produce application/json
// @Param customer body MenuItemSelection true "Customer object"
// @Success 201 {object} generated.Customer
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/customer/menu/select [post]
func SelectMenuItemHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var selection MenuItemSelection
		err := json.NewDecoder(r.Body).Decode(&selection)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Publish event to RabbitMQ
		event := broker.Event{
			Type: broker.MenuItemSelected,
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

// GetAllCustomers godoc
// @Summary Get all customers
// @Description Fetches a list of all customers from the database
// @Tags customers
// @Produce application/json
// @Success 200 {array} generated.Customer
// @Router /api/customer [get]
func GetAllCustomers(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		customers, err := queries.GetAllCustomers(ctx)
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
func GetCustomerById(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			log.Println(err)
			return
		}

		customer, err := queries.GetCustomerByID(ctx, int32(id))
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
func DeleteCustomer(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			log.Println(err)
			return
		}

		err = queries.DeleteCustomer(ctx, int32(id))
		if err != nil {
			http.Error(w, "Customer not found", http.StatusNotFound)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Customer deleted"}`))
		// w.Write(res)
	}
}

// Define individual regex patterns
var (
	minLengthRegex   = regexp.MustCompile(`^.{8,}$`)
	lowercaseRegex   = regexp.MustCompile(`[a-z]`)
	uppercaseRegex   = regexp.MustCompile(`[A-Z]`)
	numberRegex      = regexp.MustCompile(`\d`)
	specialCharRegex = regexp.MustCompile(`[@$!%*?&]`)
)

// ValidatePassword checks if the password meets all complexity requirements
func ValidatePassword(password string) error {
	if !minLengthRegex.MatchString(password) {
		return errors.New("password must be at least 8 characters long")
	}
	if !lowercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !uppercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !numberRegex.MatchString(password) {
		return errors.New("password must contain at least one number")
	}
	if !specialCharRegex.MatchString(password) {
		return errors.New("password must contain at least one special character (@$!%*?&)")
	}
	return nil
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
func CreateCustomer(queries *generated.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var customer generated.CreateCustomerParams

		err := json.NewDecoder(r.Body).Decode(&customer)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Println(err)
			return
		}

		if customer.Name == "" || customer.Email == "" || customer.Password == "" {
			http.Error(w, "All required fields must be filled", http.StatusBadRequest)
			return
		}

		// if err := ValidatePassword(customer.Password); err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		createdCustomer, err := queries.CreateCustomer(ctx, customer)
		if err != nil {
			http.Error(w, "Failed to create customer", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		subject := "Welcome to MTOGO, " + customer.Name + "!"

		body := `
    <html>
        <body style="font-family: Arial, sans-serif; color: #333;">
            <h1 style="color: #4CAF50;">Welcome to [Your Service Name]!</h1>
            <p>Hi ` + customer.Name + `,</p>
            <p>We're thrilled to have you join our community! Thank you for signing up with [Your Service Name].</p>
            
            <p>Here’s what you can look forward to as a new member:</p>
            <ul>
                <li><strong>Personalized Experience:</strong> Tailored recommendations and insights just for you.</li>
                <li><strong>Exclusive Access:</strong> Enjoy early access to new features and updates.</li>
                <li><strong>Dedicated Support:</strong> Our team is here to assist you whenever you need.</li>
            </ul>

            <p>To get started, simply log in and explore. We’re here to make sure you have a seamless experience, so don’t hesitate to reach out if you have any questions.</p>

            <p style="margin-top: 30px;">Cheers,</p>
            <p>The [Your Service Name] Team</p>
            <footer style="margin-top: 20px; font-size: 0.9em; color: #666;">
                <hr>
                <p>If you did not sign up for this account, please ignore this email.</p>
            </footer>
        </body>
    </html>
`

		err = mailer.SendMailWithGomail(customer.Email, subject, body)
		if err != nil {
			log.Println("Failed to send email:", err)
		}

		res, _ := json.Marshal(createdCustomer)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}
}
