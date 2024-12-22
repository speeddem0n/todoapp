package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt" // пакет для работы с jwt token
	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/pkg/repository"
)

const (
	salt       = "asd2gewvqwef234rtf"  // Случайная соль для пароля
	tokenTTL   = 12 * time.Hour        // Время валидности jwt токена
	signingKey = "qweg1q2egqewf#fqvcq" // Ключ для подписи jwt токена
)

type tokenClaims struct { // структура tokenClaims для послудующей передачи в NewWithClaims в методе GenerateToken
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct { // Структура AuthService в которой находится соответствующий интерфейс из репозитория
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService { // Конструктор для структуры AuthService принимает соответствующий интерфейс repository.Authorization из репозитория
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
	return token.SignedString([]byte(signingKey)) // Generate the signing string
}

func (s *AuthService) ParseToken(accessToken string) (int, error) { // Парсит jwt токен и возвращает id пользователя если все ОК
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // Проверяем что метод подписи токена HMAC
			return nil, errors.New("invalid singing method") // Ошибка неверный метод подписи
		}

		return []byte(signingKey), nil // Возвращаем ключ подписи если все ОК
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims) // Функция ParseWithClaims возвращает объект токена, в котором есть поле Claims типа interface. Приведем его в структуре tokenClaims и проверим все ли хорошо.
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string { // Функция для генерации Хэша пароля
	hash := sha1.New()
	hash.Write([]byte(password)) // Хэшируем пароль

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))) // Добовляем к хэшу пароля случайную соль и возвращаем его
}
