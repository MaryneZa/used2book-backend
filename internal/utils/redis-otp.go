package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis initializes the Redis client using environment variables.
func InitRedis() {
	// Attempt to load the .env file.
	// In production (e.g., on Kubernetes), this file might not exist,
	// so log a warning and continue.

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}


	// log.Println("ENV - redis" ,os.Getenv("ENV"))

	// if os.Getenv("ENV") != "production" {
    //     if err := godotenv.Load(); err != nil {
    //         log.Println("Warning: .env file not found, using system environment variables - redis")
    //     }
    // }

	// Retrieve Redis host and port from environment variables.
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis" // default service name in Kubernetes
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	address := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // Set this if you have a Redis password; otherwise leave blank.
		DB:       0,  // Use default DB
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", address, err)
	} else {
		fmt.Printf("🔥 Redis Connected Successfully on %s\n", address)
	}
}
