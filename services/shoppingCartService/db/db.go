package db

import (
	"github.com/go-redis/redis/v8"
)

type ShoppingCartService struct {
	redisClient *redis.Client
}

func NewShoppingCartService(client *redis.Client) *ShoppingCartService {
	return &ShoppingCartService{redisClient: client}
}
