package db

import (
	"github.com/go-redis/redis/v8"
)

type ShoppingCartRepository struct {
	redisClient *redis.Client
}

func NewShoppingCartRepository(client *redis.Client) *ShoppingCartRepository {
	return &ShoppingCartRepository{redisClient: client}
}
