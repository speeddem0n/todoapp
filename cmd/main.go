package main

import (
	"log"

	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/pkg/handler"
)

func main() {
	handlers := new(handler.Handler) // Инициализируем наш обработчик из pkg/handler
	srv := new(todo.Server)          // Инициализируеми структуру сервера
	port := "4000"
	err := srv.Run(port, handlers.InitRoutes()) // Запускаем сервер на указаном порте
	if err != nil {
		log.Fatalf("Ошибка во времая запуска серевера: %s", err.Error()) // Обработка ошибки запуска сервера
	}
	log.Printf("Сервер запущен по адресу: :%s", port)
}
