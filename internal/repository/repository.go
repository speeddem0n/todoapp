package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/speeddem0n/todoapp/internal/models"
)

const ( // Константы с названием таблиц из БД
	usersTable      = "users"
	todoListTable   = "todo_lists"
	usersListsTable = "users_list"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Authorization interface { // Интерфейс для авторизации
	CreateUser(user models.User) (int, error)               // Метод CreateUser для создания пользователя
	GetUser(username, password string) (models.User, error) // Метод GetUser для поулчения id пользоваетя
}

type TodoList interface { // Интерфейс для работы со списками
	Create(userId int, list models.TodoList) (int, error)          // Метод для создания списка возвращает id созданного списка и ошибку
	GetAll(userId int) ([]models.TodoList, error)                  // Метод для возвращения всех списков "todo" конкретного пользователя (принимает id пользователя)
	GetById(userId, listId int) (models.TodoList, error)           // Метод для получения списка пользователя по его ID
	Update(userId, listId int, input models.UpdateListInput) error // Метод для обновления списка по его id
	Delete(userId, listId int) error                               // Метод для удаления списка по его ID
}

type TodoItem interface { // Интерфейс для работы с элеметоми списков
	Create(listId int, item models.TodoItem) (int, error)          // Метод для создания "todo" элемента возвращает id созданного элемента и ошибку
	GetAll(userId, listId int) ([]models.TodoItem, error)          // Метод для возвращения всех элементов списка конкретного пользователя (принимает id пользователя и списка)
	GetById(userId, itemId int) (models.TodoItem, error)           // Метод для получения элемента списка по его ID
	Update(userId, itemId int, input models.UpdateItemInput) error // Метод для обновления элемента списка по его id
	Delete(userId, itemId int) error                               // Метод для удаления эелемнта списка по его ID
}

type Repository struct { // Структура Repository содержит 3 интерфейса аналогично структуре service (4 УРОВЕНЬ)
	Authorization // Интерфейс для авторизации
	TodoList      // Интерфейс для работы со списками
	TodoItem      // Интерфейс для работы с элеметами списков
}

func NewRepository(db *sqlx.DB) *Repository { // Конструктор для структуры Repository
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
