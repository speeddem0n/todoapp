package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// godotenv для работы с .env файлами
	_ "github.com/lib/pq"        // драйвер для работы с БД
	"github.com/sirupsen/logrus" // logrus для логирования
	config "github.com/speeddem0n/todoapp/internal/configs"
	"github.com/speeddem0n/todoapp/internal/connections"
	"github.com/speeddem0n/todoapp/internal/handler"
	"github.com/speeddem0n/todoapp/internal/repository"
	"github.com/speeddem0n/todoapp/internal/service"
	"github.com/spf13/viper" // viper для работы с config файлами
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Инициализируем конфиг
	err := config.InitConfig()
	if err != nil {
		logrus.Fatalf("error initialization config: %s", err.Error())
	}

	// Инициализируем подключение к базе данных
	db, err := connections.NewPostgresConnection()
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()

	repos := repository.NewRepository(db)    // Ицициализируем структуру БД
	services := service.NewService(repos)    // Инициализируем структуру сервисов и передаем в нее структуру БД
	handlers := handler.NewHandler(services) // Инициализируем структуру обработчиков и передаем в нее структуру сервисов

	// Инициализируем структуру сервера
	srv := http.Server{
		Addr:           ":" + viper.GetString("port"),
		Handler:        handlers.InitRoutes(),
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		logrus.Infof("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("Server stopped")
		}
	}()
	logrus.Info("TodoApp is Running")

	// Реализация Graceful shutdown
	quit := make(chan os.Signal, 1)                      // Создаем канал типа os.Signal
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT) // SIGTERM - Termination signal, SIGINT - Terminal interrupt signal
	<-quit                                               // Канал считает системный сигнал

	logrus.Print("TodoApp Shutting Down...")

	ctx, ctxCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer ctxCancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		logrus.Errorf("Error occured on server shutting down: %s", err.Error())
	}

}
