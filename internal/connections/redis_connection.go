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
	var (
		DialTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	)

	DialTimeout = (time.Second * viper.GetInt("redis.dialtimeout"))

	rdb := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.host"), // Адрес Redis
		Password:     os.Getenv("REDIS_PASSWORD"),   // Пароль, если есть
		DB:           viper.GetInt("redis.dbname"),  // Используемая база данных
		Username:     viper.GetString("redis.username"),
		MaxRetries:   viper.GetInt("redis.retries"),
		DialTimeout:  time.Second * viper.GetInt("redis.dialtimeout"),
		ReadTimeout:  viper.GetInt("redis.timeout") * time.Second(),
		WriteTimeout: time.Second * viper.GetInt("redis.timeout"),
	})

	// Проверяем соединение
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
