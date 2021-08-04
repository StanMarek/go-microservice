package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ClientRedis *redis.Client

func ConnectRedis() {
	redisdsn := os.Getenv("REDIS_DSN")
	if len(redisdsn) == 0 {
		redisdsn = "127.0.0.1:6379"
	}
	ClientRedis = redis.NewClient(&redis.Options{
		Addr:     redisdsn,
		Password: "",
		DB:       0,
	})
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc()

	if _, err := ClientRedis.Ping(ctx).Result(); err != nil {
		log.Fatal(err)
	}
}
