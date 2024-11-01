package main

import (
	"fmt"
	"log"
	"net/http"

	"soft-exam-2/db"
	_ "soft-exam-2/docs"
	"soft-exam-2/handlers"
	"github.com/rs/cors"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func run() (*http.ServeMux, error) {
	err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /api/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /api/todos", handlers.GetAllTodos)
	mux.HandleFunc("GET /api/todo/{id}", handlers.GetTodo)
	mux.HandleFunc("DELETE /api/todo/{id}", handlers.DeleteTodo)
	mux.HandleFunc("POST /api/todo", handlers.CreateTodo)

	return mux, nil
}

func main() {
	mux, err := run()
	if err != nil {
		log.Fatal(err)
	}
	// Ensure the database connection is closed at the end
	defer db.DB.Close()

	// CORS stuff
	handler := cors.Default().Handler(mux)

	fmt.Println("Running server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}