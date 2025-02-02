package connections

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func NewPostgresConnection() (*sqlx.DB, error) {

	// Считываем настройки для подключения к БД из конфига
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.username"),
		os.Getenv("DB_PASSWORD"),
		viper.GetString("db.dbname"),
		viper.GetString("db.sslmode"),
	))

	if err != nil {
		return nil, err
	}

	// Методом Ping() проверяем работоспособность подключения к БД
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
