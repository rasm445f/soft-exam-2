package db

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

func TestGetCart(t *testing.T) {
	// Create a new mock Redis client
	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Initialize the repository with the mock Redis client
	repo := &ShoppingCartRepository{redisClient: db}

	// Test data
	cart := &ShoppingCart{
		CustomerId:   123,
		RestaurantId: 456,
		TotalAmount:  58.0,
		VatAmount:    11,
		Items: []ShoppingCartItem{
			{
				Id:       1,
				Name:     "Sample Item",
				Price:    29.0,
				Quantity: 2,
			},
		},
	}

	cartData, err := json.Marshal(cart)
	if err != nil {
		t.Fatalf("unexpected error marshalling cart: %v", err)
	}

	t.Run("successful get", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectGet(cartKey).SetVal(string(cartData))

		result, err := repo.GetCart(context.Background(), 123)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if result.CustomerId != cart.CustomerId {
			t.Errorf("got CustomerId %d, want %d", result.CustomerId, cart.CustomerId)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("non-existent cart", func(t *testing.T) {
		cartKey := "cart:69"
		mock.ExpectGet(cartKey).RedisNil()

		_, err := repo.GetCart(context.Background(), 69)
		if err != redis.Nil {
			t.Error("expected redis.Nil error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
	t.Run("redis error", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectGet(cartKey).SetErr(fmt.Errorf("network error"))

		_, err := repo.GetCart(context.Background(), 123)
		if err == nil {
			t.Error("expected error, got nil")
		}

		expectedErr := "failed to retrieve shopping cart: network error"
		if err.Error() != expectedErr {
			t.Errorf("got error %v, want %v", err, expectedErr)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("unmarshal error", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectGet(cartKey).SetVal("{invalid json}")

		_, err := repo.GetCart(context.Background(), 123)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "failed to unmarshal shopping cart") {
			t.Errorf("got error %v, want error containing 'failed to unmarshal shopping cart'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}

func TestSaveCart(t *testing.T) {
	// Create a new mock Redis client
	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Initialize the repository with the mock Redis client
	repo := &ShoppingCartRepository{redisClient: db}

	// Test data
	cart := &ShoppingCart{
		CustomerId:   123,
		RestaurantId: 456,
		TotalAmount:  58.0,
		VatAmount:    11,
		Items: []ShoppingCartItem{
			{
				Id:       1,
				Name:     "Sample Item",
				Price:    29.0,
				Quantity: 2,
			},
		},
	}

	t.Run("successful save", func(t *testing.T) {
		cartData, err := json.Marshal(cart)
		if err != nil {
			t.Fatalf("unexpected error marshalling cart: %v", err)
		}

		cartKey := "cart:123"
		mock.ExpectSet(cartKey, cartData, 0).SetVal("OK")

		err = repo.SaveCart(context.Background(), cart)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("marshal error", func(t *testing.T) {
		invalidCart := &ShoppingCart{
			CustomerId: 123,
			Items: []ShoppingCartItem{
				{
					// Create an invalid field that can't be marshaled
					Price: math.Inf(1),
				},
			},
		}

		err := repo.SaveCart(context.Background(), invalidCart)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "failed to marshal cart") {
			t.Errorf("got error %v, want error containing 'failed to marshal cart'", err)
		}
	})

	t.Run("redis error", func(t *testing.T) {
		cartData, _ := json.Marshal(cart)
		cartKey := "cart:123"
		mock.ExpectSet(cartKey, cartData, 0).SetErr(fmt.Errorf("network error"))

		err := repo.SaveCart(context.Background(), cart)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}

func TestClearCart(t *testing.T) {
	// Create a new mock Redis client
	db, mock := redismock.NewClientMock()
	defer db.Close()

	// Initialize the repository with the mock Redis client
	repo := &ShoppingCartRepository{redisClient: db}

	t.Run("successful clear", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectDel(cartKey).SetVal(2)

		err := repo.ClearCart(context.Background(), 123)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("non-existent cart", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectDel(cartKey).SetVal(0)

		err := repo.ClearCart(context.Background(), 123)
		if err == nil {
			t.Error("expected error, got nil")
		}

		expectedErr := "cart for customer ID 123 does not exist"
		if err.Error() != expectedErr {
			t.Errorf("got error %v, want %v", err, expectedErr)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	// Add to TestClearCart function
	t.Run("redis error", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectDel(cartKey).SetErr(fmt.Errorf("network error"))

		err := repo.ClearCart(context.Background(), 123)
		if err == nil {
			t.Error("expected error, got nil")
		}

		expectedErr := "failed to clear cart: network error"
		if err.Error() != expectedErr {
			t.Errorf("got error %v, want %v", err, expectedErr)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}
