package main

import (
	"context"
	"fmt"
	"log"

	"github.com/oTuff/go-startkode/db"
)

func main() {
	redisClient := db.Redis_conn()
	defer redisClient.Close()

	commands := db.NewShoppingCartService(redisClient)

	// TODO: remove. Example usage
	item := db.ShoppingCartItem{
		ID:       "123",
		Name:     "Widget",
		Price:    19.99,
		Quantity: 2,
	}
	item2 := db.ShoppingCartItem{
		ID:       "1234",
		Name:     "Widgetnew",
		Price:    19.99,
		Quantity: 1,
	}

	// context is gonna come from the request: `ctx := r.Context()`
	ctx := context.Background()

	if err := commands.AddItem(ctx, "user1", item); err != nil {
		log.Fatalf("Could not add item: %v", err)
	}
	if err := commands.AddItem(ctx, "user1", item2); err != nil {
		log.Fatalf("Could not add item: %v", err)
	}

	items, err := commands.ViewCart(ctx, "user1")
	if err != nil {
		log.Fatalf("Could not retrieve cart: %v", err)
	}

	fmt.Printf("Cart items for user1: %+v\n", items)

	if err := commands.ClearCart(ctx, "user1"); err != nil {
		log.Fatalf("Could not clear cart: %v", err)
	}

	fmt.Println("Cart cleared for user1")
}
