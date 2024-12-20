package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx" // sqlx для работы с базой данных
	_ "github.com/lib/pq"     // драйвер для работы с postgreSQL
)

type Config struct { // Структура с настройками для подключения к БД
	Host     string // Адрес подключения к БД
	Port     string // Порт подключения к БД
	Username string // Имя пользователя БД
	Password string // Пароль ДБ
	DBName   string // Название базы данных
	SSLMode  string // sslmode
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) { // Функция для подключения к БД принимает Config struct
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)) // Считываем настройки для подключения к БД через fmt.Sprintf
	if err != nil {
		return nil, err // Обрабатываем ошибку
	}

	err = db.Ping() // Методом Ping() проверяем работоспособность подключения к БД
	if err != nil {
		return nil, err
	}

	return db, nil // Возвращаем *sqlx.DB
}
