package service

import (
	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/pkg/repository"
)

type Authorization interface { // интерфейс Authorization
	CreateUser(user todo.User) (int, error)                  // Создает нового пользователя и возвращает его Id
	GenerateToken(username, password string) (string, error) // Создает jwt токен
	ParseToken(token string) (int, error)                    // Парсит jwt токен и возвращает id пользователя если все ОК
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)          // Метод для создания списка возвращает id созданного списка и ошибку
	GetAll(userId int) ([]todo.TodoList, error)                  // Метод для возвращения всех списков "todo" конкретного пользователя (принимает id пользователя)
	GetById(userId, listId int) (todo.TodoList, error)           // Метод для получения списка пользователя по его ID
	Update(userId, listId int, input todo.UpdateListInput) error // Метод для обновления списка по его id
	Delete(userId, listId int) error                             // Метод для удаления списка по его ID
}

type TodoItem interface {
}

type Service struct { // Структура service содержит 3 интерфейса (3 УРОВЕНЬ)
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service { // Конструктор для структуры Service, принимает указатьль на структуры Repository что бы обратится к БД
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      newTodoListService(repos.TodoList),
	}
}
