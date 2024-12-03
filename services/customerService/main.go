package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rasm445f/soft-exam-2/broker"
	"github.com/rasm445f/soft-exam-2/db"
	"github.com/rasm445f/soft-exam-2/db/generated"
	_ "github.com/rasm445f/soft-exam-2/docs"
	"github.com/rasm445f/soft-exam-2/domain"
	"github.com/rasm445f/soft-exam-2/handlers"
	"github.com/rs/cors"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func run() (http.Handler, error) {
	db, err := db.ConnectDB()
	if err != nil {
		return nil, err
	}

	// Initialize Queries with DB
	queries := generated.New(db)
	customerDomain := domain.NewCustomerDomain(queries)
	customerHandler := handlers.NewCustomerHandler(customerDomain)

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /api/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /api/customer", customerHandler.GetAllCustomers())
	mux.HandleFunc("GET /api/customer/{id}", customerHandler.GetCustomerById())
	mux.HandleFunc("DELETE /api/customer/{id}", customerHandler.DeleteCustomer())
	mux.HandleFunc("POST /api/customer", customerHandler.CreateCustomer())
	mux.HandleFunc("PATCH /api/customer/{id}", customerHandler.UpdateCustomer())

	//CORS stuff
	handler := cors.Default().Handler(mux)

	return handler, err
}

// @title Customer Service API
// @version 1.0
// @description This is the API documentation for the Customer Service.
// @contact.name API Support
// @contact.email support@example.com
// @host localhost:8081
func main() {
	// Initialize RabbitMQ
	broker.InitRabbitMQ()
	defer broker.CloseRabbitMQ()

	mux, err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running server on port 8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
