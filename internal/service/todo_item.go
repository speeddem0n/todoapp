package service

import (
	"errors"

	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/internal/models"
	"github.com/speeddem0n/todoapp/internal/repository"
)

type todoItemService struct { // Структура TodoItemService в которой находится соответствующий интерфейс из репозитория
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func newtodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *todoItemService { // Конструктор для структуры todoItemService принимает соответствующий интерфейсы repository.TodoList и repository.TodoItem из репозитория
	return &todoItemService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (s *todoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) { // Метод для создания "todo" элемента возвращает id созданного элемента и ошибку
	_, err := s.listRepo.GetById(userId, listId) //Используем метод GetById что убедится что такоей список существует для данного пользователя
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listId, item) // Возвращает анологичный метод из репозитория

}

func (s *todoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) { // Метод для возвращения всех элементов списка конкретного пользователя (принимает id пользователя и списка)
	return s.repo.GetAll(userId, listId) // Возвращает анологичный метод из репозитория
}

func (s *todoItemService) GetById(userId, itemId int) (todo.TodoItem, error) { // Метод для получения элемента списка по его ID
	return s.repo.GetById(userId, itemId) // Возвращает анологичный метод из репозитория
}

func (s *todoItemService) Update(userId, itemId int, input models.UpdateItemInput) error { // Метод для обновления элемента списка по его id
	_, err := s.repo.GetById(userId, itemId) //Используем метод GetById что убедиться что такоей  список существует для данного пользователя
	if err != nil {
		return errors.New("list doesn't exist")
	}
	if err := input.Validate(); err != nil { // Валидация инпута
		return err
	}
	return s.repo.Update(userId, itemId, input) // Возвращает анологичный метод из репозитория
}

func (s *todoItemService) Delete(userId, itemId int) error { // Метод для удаления эелемнта списка по его ID
	return s.repo.Delete(userId, itemId)
}
