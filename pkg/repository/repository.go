package repository

import (
	"github.com/jmoiron/sqlx"
	todo "github.com/speeddem0n/todoapp"
)

const ( // Константы с названием таблиц из БД
	usersTable      = "users"
	todoListTable   = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository { // Конструктор для структуры Repository
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
