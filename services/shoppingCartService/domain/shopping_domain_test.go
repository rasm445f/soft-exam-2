package domain

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/rasm445f/soft-exam-2/db"
)

func TestAddItemDomain(t *testing.T) {
	// Create a new mock Redis client
	redisDb, mock := redismock.NewClientMock()
	defer redisDb.Close()

	// Initialize the repository with the mock Redis client
	repo := db.NewShoppingCartRepository(redisDb)
	domain := NewShoppingCartDomain(repo)

	// Test data
	cart := &db.ShoppingCart{
		CustomerId:   123,
		RestaurantId: 456,
		TotalAmount:  60.0,
		VatAmount:    12.0,
		Items: []db.ShoppingCartItem{
			{
				Id:       1,
				Name:     "Sample Item",
				Price:    30,
				Quantity: 2,
			},
		},
	}

	cartData, err := json.Marshal(cart)
	if err != nil {
		t.Fatalf("unexpected error marshalling cart: %v", err)
	}

	// add item params
	itemParams := AddItemParams{
		123,
		456,
		"Sample Item",
		30.0,
		2,
	}

	t.Run("successfully create a new cart", func(t *testing.T) {
		cartKey := "cart:123"
		// first we will get the RedisNil error because the cart does not exist
		mock.ExpectGet(cartKey).RedisNil()

		// then we should get the data
		mock.ExpectSet(cartKey, cartData, 0).SetVal("OK")

		err := domain.AddItemDomain(context.Background(), itemParams)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
	// t.Run("successfully Add item to existing cart", func(t *testing.T) {
	// 	// first create a cart like above
	// 	cartKey := "cart:123"
	// 	// first we will get the RedisNil error because the cart does not exist
	// 	mock.ExpectGet(cartKey).RedisNil()
	// 	mock.ExpectSet(cartKey, cartData, 0).SetVal("OK")
	// 	err := domain.AddItemDomain(context.Background(), itemParams)
	// 	if err != nil {
	// 		t.Errorf("unexpected error: %v", err)
	// 	}
	//
	// 	cart.TotalAmount += itemParams.Price
	// 	cart.VatAmount += itemParams.Price * 0.2
	//
	// 	cartData, err := json.Marshal(cart)
	// 	if err != nil {
	// 		t.Fatalf("unexpected error marshalling cart: %v", err)
	// 	}
	//
	// 	// insert another item in the same cart
	// 	err = domain.AddItemDomain(context.Background(), itemParams)
	// 	if err != nil {
	// 		t.Errorf("unexpected error: %v", err)
	// 	}
	//
	// })
	t.Run("check quantity", func(t *testing.T) {})
}
