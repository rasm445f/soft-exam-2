package db

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

func TestAddItem(t *testing.T) {
	// Create a new mock Redis client
	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Initialize the repository with the mock Redis client
	repo := &ShoppingCartRepository{redisClient: db}

	// Define the test item
	item := AddItemParams{
		Name:     "Sample Item",
		Price:    29,
		Quantity: 2,
		UserId:   "user123",
	}

	// Define the Redis keys based on input
	itemIdKey := "cart:user123:nextId"
	cartKey := "cart:user123"

	// Set up expected Redis behavior for Incr
	mock.ExpectIncr(itemIdKey).SetVal(1) // Expect Incr to return 1

	// Prepare the expected JSON data for HSet
	itemWithID := struct {
		Id       int64
		Name     string
		Quantity int
		Price    float64
	}{
		Id:       1,
		Name:     item.Name,
		Quantity: item.Quantity,
		Price:    item.Price,
	}

	itemData, err := json.Marshal(itemWithID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// Print the JSON data for debugging
	fmt.Printf("Expected JSON data: %s\n", itemData)

	// Set up expected Redis behavior for HSet
	mock.ExpectHSet(cartKey, "1", itemData).SetVal(1)

	// Call AddItem and check the result
	err = repo.AddItem(context.Background(), item)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func TestViewCart(t *testing.T) {
	// Create a new mock Redis client
	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Initialize the repository with the mock Redis client
	repo := &ShoppingCartRepository{redisClient: db}

	// Define the user ID and Redis cart key
	userId := "user123"
	cartKey := fmt.Sprintf("cart:%s", userId)

	// Define multiple items expected in the user's cart
	items := []ShoppingCartItem{
		{Id: "1", Name: "Item 1", Quantity: 2, Price: 29.0},
		{Id: "2", Name: "Item 2", Quantity: 1, Price: 15.5},
	}

	// Serialize each item as it would be stored in Redis
	redisData := make(map[string]string)
	for _, item := range items {
		itemData, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("unexpected error marshalling item data: %s", err)
		}
		redisData[fmt.Sprintf("%v", item.Id)] = string(itemData)
	}

	// Set up expected Redis behavior for HGetAll
	mock.ExpectHGetAll(cartKey).SetVal(redisData)

	// Call viewCart and check the result
	retrievedItems, err := repo.ViewCart(context.Background(), userId)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// Verify that the retrieved items match the expected data
	if len(retrievedItems) != len(items) {
		t.Errorf("retrieved item count does not match: got %d, want %d", len(retrievedItems), len(items))
	}
	for i, item := range items {
		if retrievedItems[i].Id != item.Id || retrievedItems[i].Name != item.Name ||
			retrievedItems[i].Quantity != item.Quantity || retrievedItems[i].Price != item.Price {
			t.Errorf("retrieved item does not match expected item at index %d: got %+v, want %+v", i, retrievedItems[i], item)
		}
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}
