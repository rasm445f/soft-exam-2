package db

import (
	"context"
	"encoding/json"
	"fmt"
)

type AddItemParams struct {
	UserId   string  `json:"user_id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func (s *ShoppingCartRepository) AddItem(ctx context.Context, item AddItemParams) error {
	// use the userid as key for future retrieval
	key := fmt.Sprintf("cart:%s", item.UserId)

	// auto increment item key - this is also saved in redis
	itemIdKey := fmt.Sprintf("cart:%s:nextId", item.UserId)
	itemId, err := s.redisClient.Incr(ctx, itemIdKey).Result()

	// new struct using newly created itemid instead of userid
	itemWithID := struct {
		Id       int64
		Name     string
		Quantity int
		Price    float64
	}{
		Id:       itemId,
		Name:     item.Name,
		Quantity: item.Quantity,
		Price:    item.Price,
	}

	itemData, err := json.Marshal(itemWithID)
	if err != nil {
		return fmt.Errorf("failed to marshal item data: %w", err)
	}

	// Store item in Redis using HSET (hash data structure)
	return s.redisClient.HSet(ctx, key, fmt.Sprintf("%d", itemId), itemData).Err()
}

func (s *ShoppingCartRepository) ViewCart(ctx context.Context, userID string) ([]ShoppingCartItem, error) {
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

func (s *ShoppingCartRepository) ClearCart(ctx context.Context, userID string) error {
	key := fmt.Sprintf("cart:%s", userID)
	return s.redisClient.Del(ctx, key).Err()
	// 	cartKey := fmt.Sprintf("cart:%s", userID)
	// itemIdKey := fmt.Sprintf("cart:%s:nextId", userID)
	//
	// // Delete the cart key and the ID counter key
	// err := s.redisClient.Del(ctx, cartKey, itemIdKey).Err()
	// if err != nil {
	// 	return fmt.Errorf("failed to clear cart: %w", err)
	// }
	//
	// return nil
}
