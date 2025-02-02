package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Функция для инициализации конфига
func initViperConfig() error {
	viper.AddConfigPath("internal/configs") // Инициализируем путь к дириктории в которой лежат config файлы
	viper.SetConfigName("config")           // Инициализируем имя config файла
	return viper.ReadInConfig()
}

func InitConfig() error {
	err := initViperConfig()
	if err != nil {
		logrus.Fatalf("error initialization config: %s", err.Error())
		return err
	}

	// Загружаем .env файл с паролем к бд
	err = godotenv.Load()
	if err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
		return err
	}

	// Задаем для логера формат json
	logrus.SetFormatter(new(logrus.JSONFormatter))
	return nil
}
