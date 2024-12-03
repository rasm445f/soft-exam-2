package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func (s *ShoppingCartRepository) GetCart(ctx context.Context, cartId int) (*ShoppingCart, error) {
	cartKey := fmt.Sprintf("cart:%d", cartId)

	cartData, err := s.redisClient.Get(ctx, cartKey).Result()
	if err == redis.Nil {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve shopping cart: %w", err)
	}

	var cart ShoppingCart
	if err := json.Unmarshal([]byte(cartData), &cart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shopping cart: %w", err) // maybe a better error message?
	}

	return &cart, nil
}

func (s *ShoppingCartRepository) SaveCart(ctx context.Context, cart *ShoppingCart) error {
	cartKey := fmt.Sprintf("cart:%d", cart.CustomerId)
	cartData, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("failed to marshal cart: %w", err)
	}
	return s.redisClient.Set(ctx, cartKey, cartData, 0).Err()
}

func (s *ShoppingCartRepository) ClearCart(ctx context.Context, customerId int) error {
	cartKey := fmt.Sprintf("cart:%d", customerId)

	// delete
	deletedCount, err := s.redisClient.Del(ctx, cartKey).Result()
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	// redisClient.Del() returns the number of deleted items
	if deletedCount == 0 {
		return fmt.Errorf("cart for customer ID %d does not exist", customerId)
	}

	return nil
}
