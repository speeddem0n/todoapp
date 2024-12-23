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
	Create(userId int, list todo.TodoList) (int, error)          // Метод для создания списка возвращает id созданного списка и ошибку
	GetAll(userId int) ([]todo.TodoList, error)                  // Метод для возвращения всех списков "todo" конкретного пользователя (принимает id пользователя)
	GetById(userId, listId int) (todo.TodoList, error)           // Метод для получения списка пользователя по его ID
	Update(userId, listId int, input todo.UpdateListInput) error // Метод для обновления списка по его id
	Delete(userId, listId int) error                             // Метод для удаления списка по его ID
}

type TodoItem interface {
	Create(listId int, item todo.TodoItem) (int, error) // Метод для создания "todo" элемента возвращает id созданного элемента и ошибку
	GetAll(userId, listId int) ([]todo.TodoItem, error) // Метод для возвращения всех элементов списка конкретного пользователя (принимает id пользователя и списка)
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
		TodoItem:      NewTodoItemPostgres(db),
	}
}
