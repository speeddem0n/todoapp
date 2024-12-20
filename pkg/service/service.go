package service

import "github.com/speeddem0n/todoapp/pkg/repository"

type Autorization interface {
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Autorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service { // Конструктор для структуры Service, принимает указатьль на структуры Repository что бы обратится к БД
	return &Service{}
}
