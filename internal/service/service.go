package service

import (
	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/internal/models"
	"github.com/speeddem0n/todoapp/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface { // Интерфейс для авторизации
	CreateUser(user todo.User) (int, error)                  // Создает нового пользователя и возвращает его Id
	GenerateToken(username, password string) (string, error) // Создает jwt токен
	ParseToken(token string) (int, error)                    // Парсит jwt токен и возвращает id пользователя если все ОК
}

type TodoList interface { // Интерфейс для работы со списками
	Create(userId int, list todo.TodoList) (int, error)            // Метод для создания списка возвращает id созданного списка и ошибку
	GetAll(userId int) ([]todo.TodoList, error)                    // Метод для возвращения всех списков "todo" конкретного пользователя (принимает id пользователя)
	GetById(userId, listId int) (todo.TodoList, error)             // Метод для получения списка пользователя по его ID
	Update(userId, listId int, input models.UpdateListInput) error // Метод для обновления списка по его id
	Delete(userId, listId int) error                               // Метод для удаления списка по его ID
}

type TodoItem interface { // Интерфейс для работы с элеметоми списков
	Create(userId, listId int, item todo.TodoItem) (int, error)    // Метод для создания "todo" элемента возвращает id созданного элемента и ошибку
	GetAll(userId, listId int) ([]todo.TodoItem, error)            // Метод для возвращения всех элементов списка конкретного пользователя (принимает id пользователя и списка)
	GetById(userId, itemId int) (todo.TodoItem, error)             // Метод для получения элемента списка по его ID
	Update(userId, itemId int, input models.UpdateItemInput) error // Метод для обновления элемента списка по его id
	Delete(userId, itemId int) error                               // Метод для удаления эелемнта списка по его ID
}

type Service struct { // Структура service содержит 3 интерфейса (3 УРОВЕНЬ)
	Authorization // Интерфейс для авторизации
	TodoList      // Интерфейс для работы со списками
	TodoItem      // Интерфейс для работы с элеметоми списков
}

func NewService(repos *repository.Repository) *Service { // Конструктор для структуры Service, принимает указатьль на структуры Repository что бы обратится к БД
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      newTodoListService(repos.TodoList),
		TodoItem:      newtodoItemService(repos.TodoItem, repos.TodoList),
	}
}
