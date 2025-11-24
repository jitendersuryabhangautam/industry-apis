package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() {
	// Remove the colon (:) to assign to the global variable, not create a local one
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // Add password if your Redis server requires it
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}
	fmt.Println("✅ Redis connected successfully")
}

func Close() {
	if Client != nil {
		_ = Client.Close()
	}
}
