package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
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

func (s *ShoppingCartRepository) UpdateCart(ctx context.Context, userID string, itemID string, newQuantity int) error {
	// The key for the user's cart
	cartKey := fmt.Sprintf("cart:%s", userID)

	// Fetch the item from the cart
	itemData, err := s.redisClient.HGet(ctx, cartKey, itemID).Result()
	if err != nil {
		// If the item does not exist
		if err == redis.Nil {
			return fmt.Errorf("item with ID %s does not exist in the cart", itemID)
		}
		return fmt.Errorf("failed to fetch item: %w", err)
	}

	// Deserialize the item JSON into a ShoppingCartItem
	var item ShoppingCartItem
	if err := json.Unmarshal([]byte(itemData), &item); err != nil {
		return fmt.Errorf("failed to unmarshal item data: %w", err)
	}

	// Update or remove the item based on newQuantity
	if newQuantity > 0 {
		// Update the item's quantity
		item.Quantity = newQuantity

		// Serialize the updated item back to JSON
		updatedItemData, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("failed to marshal updated item: %w", err)
		}

		// Store the updated item in Redis
		if err := s.redisClient.HSet(ctx, cartKey, itemID, updatedItemData).Err(); err != nil {
			return fmt.Errorf("failed to update item in cart: %w", err)
		}
	} else {
		// Remove the item from the cart if newQuantity <= 0
		if err := s.redisClient.HDel(ctx, cartKey, itemID).Err(); err != nil {
			return fmt.Errorf("failed to remove item from cart: %w", err)
		}
	}

	return nil
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
