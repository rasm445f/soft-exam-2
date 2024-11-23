package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func Redis_conn() (*redis.Client, error) {
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}

	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})
	if client == nil {
		return nil, errors.New("could not connect to client")
	}

	return client, nil
}
