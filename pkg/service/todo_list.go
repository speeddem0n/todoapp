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

func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) { // Метод для создания списка возвращает id созданного списка и ошибку
	return s.repo.Create(userId, list) // Возвращает анологичный метод из репозитория
}

func (s *TodoListService) GetAll(userId int) ([]todo.TodoList, error) { // Метод для возвращения всех списков "todo" конкретного пользователя (принимает id пользователя)
	return s.repo.GetAll(userId) // Возвращает анологичный метод из репозитория
}

func (s *TodoListService) GetById(userId, listId int) (todo.TodoList, error) { // Метод для получения списка пользователя по его ID
	return s.repo.GetById(userId, listId) // Возвращает анологичный метод из репозитория
}

func (s *TodoListService) Update(userId, listId int, input todo.UpdateListInput) error { // Метод для обновления списка по его id
	if err := input.Validate(); err != nil { // Валидация инпута
		return err
	}
	return s.repo.Update(userId, listId, input) // Возвращает анологичный метод из репозитория
}

func (s *TodoListService) Delete(userId, listId int) error { // Метод для удаления списка по его ID
	return s.repo.Delete(userId, listId) // Возвращает анологичный метод из репозитория
}
