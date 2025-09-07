package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if _, err := Rdb.Ping(Ctx).Result(); err != nil {
		fmt.Printf("❌ Failed to connect to Redis: %v\n", err)
	} else {
		fmt.Println("✅ Connected to Redis")
	}
}

func Close() {
	if Rdb != nil {
		_ = Rdb.Close()
	}
}
