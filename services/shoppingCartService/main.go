package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/oTuff/go-startkode/db"
	_ "github.com/oTuff/go-startkode/docs"
	"github.com/oTuff/go-startkode/handlers"
)

func run() (http.Handler, error) {
	redisClient := db.Redis_conn()

	commands := db.NewShoppingCartRepository(redisClient)

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /api/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("POST /api/shopping", handlers.AddItemHandler(commands))
	// TODO: implement the rest of the necessary endpoints
	// mux.HandleFunc("GET /api/todo/{id}", handlers.GetTodo(commands))
	// mux.HandleFunc("DELETE /api/todo/{id}", handlers.DeleteTodo(commands))
	// mux.HandleFunc("POST /api/todo", handlers.CreateTodo(commands))

	//CORS stuff
	handler := cors.Default().Handler(mux)

	return handler, nil // TODO: should there be a possible error value?
}

func main() {
	mux, err := run()
	if err != nil {
		log.Fatal(err)
	}
	// defer redisClient.Close()

	fmt.Println("Running server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// // TODO: remove. Example usage
// item := db.ShoppingCartItem{
// 	ID:       "123",
// 	Name:     "Widget",
// 	Price:    19.99,
// 	Quantity: 2,
// }
// item2 := db.ShoppingCartItem{
// 	ID:       "1234",
// 	Name:     "Widgetnew",
// 	Price:    19.99,
// 	Quantity: 1,
// }
//
// // context is gonna come from the request: `ctx := r.Context()`
// ctx := context.Background()
//
// if err := commands.AddItem(ctx, "user1", item); err != nil {
// 	log.Fatalf("Could not add item: %v", err)
// }
// if err := commands.AddItem(ctx, "user1", item2); err != nil {
// 	log.Fatalf("Could not add item: %v", err)
// }
//
// items, err := commands.ViewCart(ctx, "user1")
// if err != nil {
// 	log.Fatalf("Could not retrieve cart: %v", err)
// }
//
// fmt.Printf("Cart items for user1: %+v\n", items)
//
// if err := commands.ClearCart(ctx, "user1"); err != nil {
// 	log.Fatalf("Could not clear cart: %v", err)
// }
//
// fmt.Println("Cart cleared for user1")
