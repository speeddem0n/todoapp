package service

import (
	"crypto/sha1"
	"fmt"

	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/pkg/repository"
)

const salt = "asd2gewvqwef234rtf" // Случайная соль для пароля

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService { // Конструктор для структуры AuthService
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) { // Метод CreateUser()
	user.Password = s.generatePasswordHash(user.Password) // Хэшируем пароль пользователя

	return s.repo.CreateUser(user)
}

func (s *AuthService) generatePasswordHash(password string) string { // Метод для генерации Хэша пароля
	hash := sha1.New()
	hash.Write([]byte(password)) // Хэшируем пароль

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))) // Добовляем к хэшу пароля случайную соль и возвращаем его
}
