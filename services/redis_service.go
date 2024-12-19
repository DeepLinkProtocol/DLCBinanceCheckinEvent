package services

import (
	"context"
	"fmt"
	"time"

	"go-signin-service/config"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// GetRedisIntValue retrieves a value from Redis
func GetRedisIntValue(key string) (int, error) {
	client := config.GetRedisClient()
	return client.Get(ctx, key).Int()
}

// SetRedisValue sets a key-value pair in Redis with an optional expiration time
func SetRedisValue(key string, value string, expiration time.Duration) error {
	client := config.GetRedisClient()
	return client.Set(ctx, key, value, expiration).Err()
}

func IncrRedisValue(key string) error {
	client := config.GetRedisClient()
	_, err := client.Incr(context.Background(), key).Result()
	return err
}

func GetSignInDateKey(walletAddress string) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return fmt.Sprintf("signIn:%s:%s", walletAddress, today.Format("20060102"))
}

func GetSignInCountKey(walletAddress string) string {
	return fmt.Sprintf("signInCount:%s", walletAddress)
}

func SignedToday(walletAddress string) (bool, error) {
	client := config.GetRedisClient()
	key := GetSignInDateKey(walletAddress)
	_, err := client.Get(ctx, key).Int()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func SignIn(walletAddress string) error {
	client := config.GetRedisClient()
	client.Set(ctx, GetSignInDateKey(walletAddress), 1, 0)
	signInCountKey := GetSignInCountKey(walletAddress)
	return client.Incr(ctx, signInCountKey).Err()
}
