package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitRedisClient initializes the Redis client
func InitRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",   // Update to your Redis server address
		Password: "Deeplink@2024!@#", // Redis password, if any
		DB:       0,                  // Default DB
	})
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
}

// GetRedisClient returns the initialized Redis client
func GetRedisClient() *redis.Client {
	return redisClient
}
