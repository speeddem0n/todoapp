package connections

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// Инициализация Redis
func NewRedisConnection() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.host"), // Адрес Redis
		Password:     os.Getenv("REDIS_PASSWORD"),   // Пароль, если есть
		DB:           viper.GetInt("redis.dbname"),  // Используемая база данных
		Username:     viper.GetString("redis.username"),
		MaxRetries:   viper.GetInt("redis.retries"),
		DialTimeout:  time.Second * viper.GetInt("redis.dialtimeout"),
		ReadTimeout:  time.Second * viper.GetInt("redis.timeout"),
		WriteTimeout: time.Second * viper.GetInt("redis.timeout"),
	})

	// Проверяем соединение
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
