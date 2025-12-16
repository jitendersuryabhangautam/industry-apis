// Package cache provides Redis caching functionality for the application.
// It handles cache initialization, connectivity, and provides a global client instance.
package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client is the global Redis client instance used throughout the application for caching.
// It is initialized in the Init() function and should be closed using Close() on shutdown.
var Client *redis.Client

// Init initializes the Redis client by reading the connection address from environment variables.
// It creates a new Redis client, attempts to connect, and verifies connectivity with a ping.
// If connection fails, the function will panic with an error message.
func Init() {
	// Create a new Redis client with configuration from environment variables
	// REDIS_ADDR should be in the format "hostname:port" (e.g., "localhost:6379")
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // Add password if your Redis server requires authentication
		DB:       0,  // Default database index
	})

	// Create a context with a 2-second timeout for the ping operation
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Verify Redis connectivity by sending a ping command
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		// Panic if connection fails - this indicates a critical startup error
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}
	fmt.Println("✅ Redis connected successfully")
}

// Close gracefully closes the Redis client connection.
// This should be called when the application is shutting down to release Redis resources.
func Close() {
	if Client != nil {
		// Close the Redis connection (error is ignored as per original implementation)
		_ = Client.Close()
	}
}
