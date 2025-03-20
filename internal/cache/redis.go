package cache

import (
	"context"
	"os"
	"strconv"
	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	pswd := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pswd,
		DB:       db,
	})
}
