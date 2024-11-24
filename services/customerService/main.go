package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oTuff/go-startkode/db"
	"github.com/oTuff/go-startkode/db/generated"
	_ "github.com/oTuff/go-startkode/docs"
	"github.com/oTuff/go-startkode/domain"
	"github.com/oTuff/go-startkode/handlers"
	"github.com/rasm445f/soft-exam-2/broker"
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

func main() {
	// Initialize RabbitMQ
	broker.InitRabbitMQ()
	defer broker.CloseRabbitMQ()
	
	mux, err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
