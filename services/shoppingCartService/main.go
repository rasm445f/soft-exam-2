package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

// TODO: Refactor into components
var ctx = context.Background()

type ShoppingCartItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type ShoppingCartService struct {
	redisClient *redis.Client
}

func NewShoppingCartService() *ShoppingCartService {
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisHost := os.Getenv("REDIS_HOST") // Default is localhost or use your Docker hostname
	redisPort := os.Getenv("REDIS_PORT") // Default is 6379

	// TODO: simplify this
	if redisHost == "" {
		redisHost = "localhost" // Or "redis" if your Redis container is named "redis" in docker-compose
	}
	if redisPort == "" {
		redisPort = "6379"
	}

	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort), // Redis server address
		Password: redisPassword,                              // Password if set
		DB:       0,                                          // Default DB
	})
	return &ShoppingCartService{redisClient: client}
}

func (s *ShoppingCartService) AddItem(userID string, item ShoppingCartItem) error {
	key := fmt.Sprintf("cart:%s", userID)

	// Convert item to JSON
	itemData, err := json.Marshal(item)
	if err != nil {
		return err
	}

	// Store item in Redis using HSET (hash data structure)
	return s.redisClient.HSet(ctx, key, item.ID, itemData).Err()
}

func (s *ShoppingCartService) ViewCart(userID string) ([]ShoppingCartItem, error) {
	key := fmt.Sprintf("cart:%s", userID)
	itemsData, err := s.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// Deserialize JSON items into ShoppingCartItem structs
	var items []ShoppingCartItem
	for _, itemJSON := range itemsData {
		var item ShoppingCartItem
		if err := json.Unmarshal([]byte(itemJSON), &item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *ShoppingCartService) ClearCart(userID string) error {
	key := fmt.Sprintf("cart:%s", userID)
	return s.redisClient.Del(ctx, key).Err()
}
func main() {
	service := NewShoppingCartService()

	// Example usage
	item := ShoppingCartItem{
		ID:       "123",
		Name:     "Widget",
		Price:    19.99,
		Quantity: 2,
	}
	item2 := ShoppingCartItem{
		ID:       "1234",
		Name:     "Widgetnew",
		Price:    19.99,
		Quantity: 1,
	}

	if err := service.AddItem("user1", item); err != nil {
		log.Fatalf("Could not add item: %v", err)
	}
	if err := service.AddItem("user1", item2); err != nil {
		log.Fatalf("Could not add item: %v", err)
	}

	items, err := service.ViewCart("user1")
	if err != nil {
		log.Fatalf("Could not retrieve cart: %v", err)
	}

	fmt.Printf("Cart items for user1: %+v\n", items)

	// if err := service.ClearCart("user1"); err != nil {
	// 	log.Fatalf("Could not clear cart: %v", err)
	// }
	//
	// fmt.Println("Cart cleared for user1")
}
