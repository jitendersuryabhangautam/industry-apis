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
// If connection fails, it logs a warning but does not panic, allowing the application to continue without caching.
func Init() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		fmt.Println("⚠️  REDIS_ADDR not configured; skipping Redis initialization")
		return
	}

	// Create a new Redis client with configuration from environment variables
	// REDIS_ADDR should be in the format "hostname:port" (e.g., "localhost:6379")
	Client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // Add password if your Redis server requires authentication
		DB:       0,  // Default database index
	})

	// Create a context with a 2-second timeout for the ping operation
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Verify Redis connectivity by sending a ping command
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		// Log warning if connection fails, but continue without caching
		fmt.Printf("⚠️  Failed to connect to Redis: %v - Application will continue without caching\n", err)
		Client = nil // Set Client to nil so handlers can check if Redis is available
		return
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
