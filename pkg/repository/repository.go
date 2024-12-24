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

type Authorization interface { // Интерфейс для авторизации
	CreateUser(user todo.User) (int, error)               // Метод CreateUser для создания пользователя
	GetUser(username, password string) (todo.User, error) // Метод GetUser для поулчения id пользоваетя
}

type TodoList interface { // Интерфейс для работы со списками
	Create(userId int, list todo.TodoList) (int, error)          // Метод для создания списка возвращает id созданного списка и ошибку
	GetAll(userId int) ([]todo.TodoList, error)                  // Метод для возвращения всех списков "todo" конкретного пользователя (принимает id пользователя)
	GetById(userId, listId int) (todo.TodoList, error)           // Метод для получения списка пользователя по его ID
	Update(userId, listId int, input todo.UpdateListInput) error // Метод для обновления списка по его id
	Delete(userId, listId int) error                             // Метод для удаления списка по его ID
}

type TodoItem interface { // Интерфейс для работы с элеметоми списков
	Create(listId int, item todo.TodoItem) (int, error) // Метод для создания "todo" элемента возвращает id созданного элемента и ошибку
	GetAll(userId, listId int) ([]todo.TodoItem, error) // Метод для возвращения всех элементов списка конкретного пользователя (принимает id пользователя и списка)
	GetById(userId, itemId int) (todo.TodoItem, error)  // Метод для получения элемента списка по его ID
	Delete(userId, itemId int) error                    // Метод для удаления эелемнта списка по его ID
}

type Repository struct { // Структура Repository содержит 3 интерфейса аналогичн структуре service (4 УРОВЕНЬ)
	Authorization // Интерфейс для авторизации
	TodoList      // Интерфейс для работы со списками
	TodoItem      // Интерфейс для работы с элеметоми списков
}

func NewRepository(db *sqlx.DB) *Repository { // Конструктор для структуры Repository
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
