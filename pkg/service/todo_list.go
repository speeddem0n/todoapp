package service

import (
	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/pkg/repository"
)

type TodoListService struct { // Структура TodoListService в которой находится соответствующий интерфейс из репозитория
	repo repository.TodoList
}

func newTodoListService(repo repository.TodoList) *TodoListService { // Конструктор для структуры AuthService принимает соответствующий интерфейс repository.TodoList из репозитория
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) { // Метод для создания списка todo передает данные на след уровень, в репозиторий
	return s.repo.Create(userId, list) // Возвращает анологичный метод из репозитория
}
