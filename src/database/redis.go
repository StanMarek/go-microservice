package database

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

// exported variables must have name starting with Big letter
// other way it won't be expoerted
var ClientRedis *redis.Client
var CtxRedis = context.Background()

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
	if _, err := ClientRedis.Ping(CtxRedis).Result(); err != nil {
		log.Fatal(err)
	}
}
