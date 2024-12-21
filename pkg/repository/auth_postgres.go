package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	todo "github.com/speeddem0n/todoapp"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todo.User) (int, error) { // Метод для создания нового пользователя и записи его данных в БД
	var id int                                                                                                         // Переменная для последующей записи id пользователя
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES($1, $2, $3) RETURNING id", usersTable) // SQL запрос для добовления новго пользователя в базу данных в таблицу users
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)                                               // Выполняется SQL запрос и возвращает id только что добавленной записи в переменную row *sql.Row
	err := row.Scan(&id)                                                                                               // Достаем id юзера из row в переменную id
	if err != nil {
		return 0, err
	}
	return id, nil
}
