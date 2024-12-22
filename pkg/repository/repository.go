package repository

import (
	"github.com/jmoiron/sqlx"
	todo "github.com/speeddem0n/todoapp"
)

const ( // Константы с названием таблиц из БД
	usersTable      = "users"
	todoListTable   = "todo_lists"
	usersListsTable = "users_list"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)               // Метод CreateUser для создания пользователя
	GetUser(username, password string) (todo.User, error) // Метод GetUser для поулчения id пользоваетя
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error) // Метод для создания списка возвращает id созданного списка и ошибку
	GetAll(userId int) ([]todo.TodoList, error)         // Метод для возвращения всех списков дел конкретного пользователя (принимает id пользователя)
}

type TodoItem interface {
}

type Repository struct { // Структура Repository содержит 3 интерфейса аналогичн структуре service (4 УРОВЕНЬ)
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository { // Конструктор для структуры Repository
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
	}
}
