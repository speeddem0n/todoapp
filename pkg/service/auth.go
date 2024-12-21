package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt" // пакет для работы с jwt token
	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/pkg/repository"
)

const (
	salt       = "asd2gewvqwef234rtf" // Случайная соль для пароля
	tokenTTL   = 12 * time.Hour
	signingKey = "qweg1q2egqewf#fqvcq"
)

type tokenClaims struct { // структура tokenClaims для послудующей передачи в NewWithClaims в методе GenerateToken
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService { // Конструктор для структуры AuthService
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) { // Метод CreateUser()
	user.Password = generatePasswordHash(user.Password) // Хэшируем пароль пользователя

	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) { // Метод для генерации JWT токена
	user, err := s.repo.GetUser(username, generatePasswordHash(password)) // Используем метод Get user из repo
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), // токен перестанет быть валидным через 12 часов
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString(signingKey) // Generate the signing string
}

func generatePasswordHash(password string) string { // Функция для генерации Хэша пароля
	hash := sha1.New()
	hash.Write([]byte(password)) // Хэшируем пароль

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))) // Добовляем к хэшу пароля случайную соль и возвращаем его
}
