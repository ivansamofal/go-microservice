package cache

import (
	"context"
	"os"
	"strconv"
	"github.com/go-redis/redis/v8"
)

var (
	// RedisClient – глобальный клиент для работы с Redis.
	RedisClient *redis.Client
	Ctx         = context.Background()
)


// InitRedis инициализирует подключение к Redis.
func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	pswd := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0 // или можно обработать ошибку и установить значение по умолчанию
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr, // адрес вашего Redis сервера
		Password: pswd,               // если пароль не установлен, оставьте пустым
		DB:       db,                // используемая база
	})
}
