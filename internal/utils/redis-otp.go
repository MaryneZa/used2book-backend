package utils

import (
	"context"
	"fmt"
	"log"
	"github.com/go-redis/redis/v8"

)

var RedisClient *redis.Client
var ctx = context.Background()

// ✅ Initialize Redis
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Default Redis port
		Password: "",               // No password by default
		DB:       0,                 // Use default DB
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	} else {
		fmt.Println("🔥 Redis Connected Successfully")
	}
}
