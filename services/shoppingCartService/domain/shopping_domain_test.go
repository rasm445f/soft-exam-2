package domain

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-redis/redis/v8"
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
				Price:    30.0,
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
	t.Run("successfully Add item to existing cart", func(t *testing.T) {
		// first create a cart like above
		cartKey := "cart:123"
		// first we will get the RedisNil error because the cart does not exist
		mock.ExpectGet(cartKey).RedisNil()
		mock.ExpectSet(cartKey, cartData, 0).SetVal("OK")
		err := domain.AddItemDomain(context.Background(), itemParams)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// the expected new cart data
		newItem := db.ShoppingCartItem{
			Id:       len(cart.Items) + 1,
			Name:     itemParams.Name,
			Price:    itemParams.Price,
			Quantity: itemParams.Quantity,
		}
		cart.Items = append(cart.Items, newItem)
		cart.TotalAmount += newItem.Price * float64(newItem.Quantity)
		cart.VatAmount += newItem.Price * float64(newItem.Quantity) * 0.2

		newCartData, err := json.Marshal(cart)
		if err != nil {
			t.Fatalf("unexpected error marshalling cart: %v", err)
		}

		// calling ExpectGet should return the current cartData
		mock.ExpectGet(cartKey).SetVal(string(cartData))
		// When calling  adding another domain newCartData is expected
		mock.ExpectSet(cartKey, newCartData, 0).SetVal("OK")
		// insert another item in the same cart
		err = domain.AddItemDomain(context.Background(), itemParams)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("check quantity", func(t *testing.T) {
		itemParams.Quantity = 0

		// should return error when quantity is set to 0
		err := domain.AddItemDomain(context.Background(), itemParams)
		if err == nil {
			t.Errorf("expected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}

func TestUpdateCartDomain(t *testing.T) {
	redisDb, mock := redismock.NewClientMock()
	defer redisDb.Close()

	repo := db.NewShoppingCartRepository(redisDb)
	domain := NewShoppingCartDomain(repo)

	cart := &db.ShoppingCart{
		CustomerId:   123,
		RestaurantId: 456,
		Items: []db.ShoppingCartItem{
			{
				Id:       1,
				Name:     "Sample Item",
				Price:    30.0,
				Quantity: 2,
			},
		},
		TotalAmount: 60.0,
		VatAmount:   12.0,
	}

	cartData, err := json.Marshal(cart)
	if err != nil {
		t.Fatalf("unexpected error marshalling cart: %v", err)
	}

	t.Run("update item quantity", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectGet(cartKey).SetVal(string(cartData))

		// Modify item quantity
		updatedQuantity := 3
		cart.Items[0].Quantity = updatedQuantity
		cart.TotalAmount = cart.Items[0].Price * float64(updatedQuantity)
		cart.VatAmount = cart.TotalAmount * 0.20

		updatedCartData, err := json.Marshal(cart)
		if err != nil {
			t.Fatalf("unexpected error marshalling cart: %v", err)
		}

		mock.ExpectSet(cartKey, updatedCartData, 0).SetVal("OK")

		err = domain.UpdateCartDomain(context.Background(), 123, 1, updatedQuantity)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("negative quantity", func(t *testing.T) {
		got := domain.UpdateCartDomain(context.Background(), 123, 1, -1)
		want := "quantity cannot be negative"
		if got.Error() != want {
			t.Errorf("want '%v' got '%v'", want, got)
		}
	})

	t.Run("non-existing itemId", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectGet(cartKey).SetVal(string(cartData))

		got := domain.UpdateCartDomain(context.Background(), 123, 69, 2)
		want := "item not found in cart"
		if got.Error() != want {
			t.Errorf("want '%v' got '%v'", want, got)
		}
	})

	t.Run("remove item from cart", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectGet(cartKey).SetVal(string(cartData))

		// Remove item
		cart.Items = []db.ShoppingCartItem{}
		cart.TotalAmount = 0
		cart.VatAmount = 0

		updatedCartData, err := json.Marshal(cart)
		if err != nil {
			t.Fatalf("unexpected error marshalling cart: %v", err)
		}

		mock.ExpectSet(cartKey, updatedCartData, 0).SetVal("OK")

		err = domain.UpdateCartDomain(context.Background(), 123, 1, 0)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}

func TestViewCartDomain(t *testing.T) {
	redisDb, mock := redismock.NewClientMock()
	defer redisDb.Close()

	repo := db.NewShoppingCartRepository(redisDb)
	domain := NewShoppingCartDomain(repo)

	cart := &db.ShoppingCart{
		CustomerId:   123,
		RestaurantId: 456,
		Items: []db.ShoppingCartItem{
			{
				Id:       1,
				Name:     "Sample Item",
				Price:    30.0,
				Quantity: 2,
			},
		},
		TotalAmount: 60.0,
		VatAmount:   12.0,
	}

	cartData, err := json.Marshal(cart)
	if err != nil {
		t.Fatalf("unexpected error marshalling cart: %v", err)
	}

	t.Run("view existing cart", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectGet(cartKey).SetVal(string(cartData))

		result, err := domain.ViewCartDomain(context.Background(), 123)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if result.CustomerId != cart.CustomerId || result.TotalAmount != cart.TotalAmount {
			t.Errorf("expected %v, got %v", cart, result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("view non-existing cart", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectGet(cartKey).RedisNil()

		_, err := domain.ViewCartDomain(context.Background(), 123)
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !errors.Is(err, redis.Nil) {
			t.Errorf("expected redis.Nil, got %v", err)
		}
	})
}

func TestClearCartDomain(t *testing.T) {
	redisDb, mock := redismock.NewClientMock()
	defer redisDb.Close()

	repo := db.NewShoppingCartRepository(redisDb)
	domain := NewShoppingCartDomain(repo)

	t.Run("successfully clear cart", func(t *testing.T) {
		cartKey := "cart:123"
		mock.ExpectDel(cartKey).SetVal(1)

		err := domain.ClearCartDomain(context.Background(), 123)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}
