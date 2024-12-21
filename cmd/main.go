package main

import (
	"os"

	"github.com/joho/godotenv"   // godotenv для работы с .env файлами
	_ "github.com/lib/pq"        // драйвер для работы с БД
	"github.com/sirupsen/logrus" // Сторонний логер для логирования
	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/pkg/handler"
	"github.com/speeddem0n/todoapp/pkg/repository"
	"github.com/speeddem0n/todoapp/pkg/service"
	"github.com/spf13/viper" // viper для работы с config файлами
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter)) // Задаем для логера формат json

	err := intConfig() // Инициализируем конфиг функцией intConfig()
	if err != nil {
		logrus.Fatalf("error initialization config: %s", err.Error())
	}

	err = godotenv.Load() // Загружаем .env файл с паролем к бд
	if err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{ // Инициализируем новое подключение к базе данных и передаем в него параметры из конфига с помощью viper
		Host:     viper.GetString("db.host"),     // Адрес хост БД
		Port:     viper.GetString("db.port"),     // Порт БД
		Username: viper.GetString("db.username"), // username БД
		DBName:   viper.GetString("db.dbname"),   // Название БД
		SSLMode:  viper.GetString("db.sslmode"),  // SSLmode
		Password: os.Getenv("DB_PASSWORD"),       // Передаем пароль из .env с помощью godotenv
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)    // Ицициализируем структуру БД (4 УРОВЕНЬ)
	services := service.NewService(repos)    // Инициализируем структуру сервисов и передаем в нее структуру БД (3 УРОВЕНЬ)
	handlers := handler.NewHandler(services) // Инициализируем структуру обработчиков и передаем в нее структуру сервисов (2 УРОВЕНЬ)

	srv := new(todo.Server)                                       // Инициализируеми структуру сервера
	err = srv.Run(viper.GetString("port"), handlers.InitRoutes()) // viper.GetString(port) получает значения port из config. Запускаем сервер на указаном порте.
	if err != nil {
		logrus.Fatalf("Ошибка во времая запуска серевера: %s", err.Error()) // Обработка ошибки запуска сервера
	}

}

func intConfig() error { // Функция для инициализации конфига
	viper.AddConfigPath("configs") // Инициалицируем путь к дириктории в которой лежат config файлы
	viper.SetConfigName("config")  // Инициалицируем имя config файла
	return viper.ReadInConfig()
}
