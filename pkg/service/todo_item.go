package service

import (
	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/pkg/repository"
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

	return s.repo.Create(listId, item)

}

func (s *todoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) { // Метод для возвращения всех элементов списка конкретного пользователя (принимает id пользователя и списка)
	return s.repo.GetAll(userId, listId)
}
