package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rasm445f/soft-exam-2/db"
	"github.com/rasm445f/soft-exam-2/db/generated"
	_ "github.com/rasm445f/soft-exam-2/docs"
	"github.com/rasm445f/soft-exam-2/handlers"
	"github.com/rasm445f/soft-exam-2/domain"
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
		orderDomain := domain.NewOrderDomain(queries)
		orderHandler := handlers.NewOrderHandler(orderDomain)

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /api/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /api/order/consume", orderHandler.ConsumeOrder())
	//mux.HandleFunc("GET /api/todo/{id}", handlers.GetTodo(queries))
	//mux.HandleFunc("DELETE /api/todo/{id}", handlers.DeleteTodo(queries))
	//mux.HandleFunc("POST /api/todo", handlers.CreateTodo(queries))

	//CORS stuff
	handler := cors.Default().Handler(mux)

	return handler, err
}

func main() {
	mux, err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
