package db

import (
	"context"
	"encoding/json"
	"fmt"
)

func (s *ShoppingCartService) AddItem(ctx context.Context, userID string, item ShoppingCartItem) error {
	key := fmt.Sprintf("cart:%s", userID)

	// Convert item to JSON
	itemData, err := json.Marshal(item)
	if err != nil {
		return err
	}

	// Store item in Redis using HSET (hash data structure)
	return s.redisClient.HSet(ctx, key, item.ID, itemData).Err()
}
func (s *ShoppingCartService) ViewCart(ctx context.Context, userID string) ([]ShoppingCartItem, error) {
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

func (s *ShoppingCartService) ClearCart(ctx context.Context, userID string) error {
	key := fmt.Sprintf("cart:%s", userID)
	return s.redisClient.Del(ctx, key).Err()
}
