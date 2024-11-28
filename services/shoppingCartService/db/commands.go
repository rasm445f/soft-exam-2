package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type AddItemParams struct {
	CustomerId   int     `json:"customerId"`
	RestaurantId int     `json:"restaurantId"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
}

func (s *ShoppingCartRepository) AddItem(ctx context.Context, params AddItemParams) error {
	cartKey := fmt.Sprintf("cart:%d", params.CustomerId)

	cartData, err := s.redisClient.Get(ctx, cartKey).Result()
	var cart ShoppingCart

	if err == redis.Nil {
		// If no cart exists, create a new one
		cart = ShoppingCart{
			CustomerId:   params.CustomerId,
			RestaurantId: params.RestaurantId,
			Items:        []ShoppingCartItem{},
		}
	} else if err != nil {
		return fmt.Errorf("failed to retrieve shopping cart: %w", err)
	} else {
		// If cart exists, unmarshal the data
		if err := json.Unmarshal([]byte(cartData), &cart); err != nil {
			return fmt.Errorf("failed to unmarshal shopping cart: %w", err)
		}

		// Remove or refine the restaurant consistency logic
		// You could log a warning if restaurants differ but still proceed:
		if cart.RestaurantId != params.RestaurantId {
			// Log the inconsistency but allow adding the item
			fmt.Printf("Warning: cart belongs to restaurant %d, adding item from restaurant %d\n", cart.RestaurantId, params.RestaurantId)
		}
	}

	// Generate a new item ID for this customer's cart
	itemIdKey := fmt.Sprintf("cart:%d:nextId", params.CustomerId)
	itemId, err := s.redisClient.Incr(ctx, itemIdKey).Result()
	if err != nil {
		return fmt.Errorf("failed to increment item ID: %w", err)
	}

	// Add the new item to the cart
	newItem := ShoppingCartItem{
		Id:       int(itemId),
		Name:     params.Name,
		Price:    params.Price,
		Quantity: params.Quantity,
	}
	cart.Items = append(cart.Items, newItem)

	// Recalculate totals
	// cart.TotalAmount = 0.0
	for _, item := range cart.Items {
		cart.TotalAmount += item.Price * float64(item.Quantity)
	}
	cart.VatAmount = int(float64(cart.TotalAmount) * 0.20) // Adjust VAT calculation as needed

	// Marshal the updated cart back into JSON
	updatedCartData, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("failed to marshal updated shopping cart: %w", err)
	}

	// Save the updated cart back to Redis
	if err := s.redisClient.Set(ctx, cartKey, updatedCartData, 0).Err(); err != nil {
		return fmt.Errorf("failed to save shopping cart: %w", err)
	}

	return nil
}

func (s *ShoppingCartRepository) ViewCart(ctx context.Context, customerId int) (*ShoppingCart, error) {
	key := fmt.Sprintf("cart:%d", customerId)
	shoppingCartData, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// Deserialize JSON items into ShoppingCartItem structs
	var shoppingCart ShoppingCart
	if err := json.Unmarshal([]byte(shoppingCartData), &shoppingCart); err != nil {
		return nil, err
	}

	return &shoppingCart, nil
}

func (s *ShoppingCartRepository) UpdateCart(ctx context.Context, customerId int, itemID int, newQuantity int) error {
	// Cart key based on the customerId
	cartKey := fmt.Sprintf("cart:%d", customerId)

	// Use WATCH for atomic operations
	err := s.redisClient.Watch(ctx, func(tx *redis.Tx) error {
		// Fetch the cart data
		cartData, err := tx.Get(ctx, cartKey).Result()
		if err == redis.Nil {
			return fmt.Errorf("cart does not exist for customer ID %d", customerId)
		} else if err != nil {
			return fmt.Errorf("failed to retrieve shopping cart: %w", err)
		}

		// Deserialize the cart JSON
		var cart ShoppingCart
		if err := json.Unmarshal([]byte(cartData), &cart); err != nil {
			return fmt.Errorf("failed to unmarshal shopping cart: %w", err)
		}

		// Find the item in the cart
		itemFound := false
		for i, item := range cart.Items {
			if item.Id == itemID {
				itemFound = true
				if newQuantity > 0 {
					// Update the item's quantity
					cart.Items[i].Quantity = newQuantity
				} else {
					// Remove the item if quantity is zero or less
					cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
				}
				break
			}
		}

		if !itemFound {
			return fmt.Errorf("item with ID %d does not exist in the cart", itemID)
		}

		// Recalculate totals
		cart.TotalAmount = 0
		for _, item := range cart.Items {
			cart.TotalAmount += item.Price * float64(item.Quantity)
		}
		cart.VatAmount = int(float64(cart.TotalAmount) * 0.20) // Adjust VAT calculation

		// Serialize the updated cart
		updatedCartData, err := json.Marshal(cart)
		if err != nil {
			return fmt.Errorf("failed to marshal updated shopping cart: %w", err)
		}

		// Save the updated cart atomically
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, cartKey, updatedCartData, 0)
			return nil
		})
		return err
	}, cartKey)

	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	return nil
}

func (s *ShoppingCartRepository) ClearCart(ctx context.Context, customerId int) error {
	cartKey := fmt.Sprintf("cart:%d", customerId)
	nextIdKey := fmt.Sprintf("cart:%d:nextId", customerId)

	deletedCount, err := s.redisClient.Del(ctx, cartKey, nextIdKey).Result()
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	if deletedCount == 0 {
		return fmt.Errorf("cart for customer ID %d does not exist", customerId)
	}

	return nil
}
