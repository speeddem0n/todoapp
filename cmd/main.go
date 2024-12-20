package main

import (
	"log"

	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/pkg/handler"
	"github.com/speeddem0n/todoapp/pkg/repository"
	"github.com/speeddem0n/todoapp/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	err := intConfig()
	if err != nil {
		log.Fatalf("error initialization config: %s", err.Error())
	}
	repos := repository.NewRepository()      // Ицициализируем структуру БД
	services := service.NewService(repos)    // Инициализируем структуру сервисов и передаем в нее структуру БД
	handlers := handler.NewHandler(services) // Инициализируем структуру обработчиков и передаем в нее структуру сервисов

	srv := new(todo.Server) // Инициализируеми структуру сервера
	port := "4000"
	err = srv.Run(viper.GetString(port), handlers.InitRoutes()) // viper.GetString(port) получает значения port из config. Запускаем сервер на указаном порте.
	if err != nil {
		log.Fatalf("Ошибка во времая запуска серевера: %s", err.Error()) // Обработка ошибки запуска сервера
	}
	log.Printf("Сервер запущен по адресу: :%s", port)
}

func intConfig() error { // Функция для инициализации конфига
	viper.AddConfigPath("configs") // Инициалицируем путь к дириктории в которой лежат config файлы
	viper.SetConfigName("config")  // Инициалицируем имя config файла
	return viper.ReadInConfig()
}
