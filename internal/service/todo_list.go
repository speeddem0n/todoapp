package service

import (
	"errors"

	"github.com/speeddem0n/todoapp/internal/models"
	"github.com/speeddem0n/todoapp/internal/repository"
)

type TodoListService struct { // Структура TodoListService в которой находится соответствующий интерфейс из репозитория
	repo repository.TodoList
}

func newTodoListService(repo repository.TodoList) *TodoListService { // Конструктор для структуры TodoListService принимает соответствующий интерфейс repository.TodoList из репозитория
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list models.CreateListInput) (int, error) { // Метод для создания списка возвращает id созданного списка и ошибку
	return s.repo.Create(userId, list) // Возвращает анологичный метод из репозитория
}

func (s *TodoListService) GetAll(userId int) ([]models.TodoList, error) { // Метод для возвращения всех списков "todo" конкретного пользователя (принимает id пользователя)
	return s.repo.GetAll(userId) // Возвращает анологичный метод из репозитория
}

func (s *TodoListService) GetById(userId, listId int) (models.TodoList, error) { // Метод для получения списка пользователя по его ID
	return s.repo.GetById(userId, listId) // Возвращает анологичный метод из репозитория
}

func (s *TodoListService) Update(userId, listId int, input models.UpdateListInput) error { // Метод для обновления списка по его id
	_, err := s.repo.GetById(userId, listId) //Используем метод GerById что убедится что такоей список существует для данного пользователя
	if err != nil {
		return errors.New("list doesn't exist")
	}
	if err := input.Validate(); err != nil { // Валидация инпута
		return err
	}
	return s.repo.Update(userId, listId, input) // Возвращает анологичный метод из репозитория
}

func (s *TodoListService) Delete(userId, listId int) error { // Метод для удаления списка по его ID
	return s.repo.Delete(userId, listId) // Возвращает анологичный метод из репозитория
}
