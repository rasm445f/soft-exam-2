package db

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func Redis_conn() *redis.Client {
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	// TODO: simplify this
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

	// maybe using url can simplify connection stuff
	//  url := "redis://user:password@localhost:6379/0?protocol=3"
	// opts, err := redis.ParseURL(url)
	// if err != nil {
	//     panic(err)
	// }
	//
	// return redis.NewClient(opts)
	return client
}
