package db

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Connect to Redis
func ConnectRedis(ctx context.Context, conn string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		// Addr: "localhost:6379",
		Addr: conn,
	})

	ctxConn, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := client.Ping(ctxConn).Result()
	if err != nil {
		log.Fatal(err)
	}

	return client
}
