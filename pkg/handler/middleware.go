package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) { // Метод идентификации пользователя по jwt токену
	header := c.GetHeader(authorizationHeader) // Поулчаем header
	if header == "" {                          // Проверяем что хедер не пустой
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")                // Разделяем строку header по пробелам
	if len(headerParts) != 2 || headerParts[0] != "Bearer" { // При корректном формате хедера strings.Split должен вернуть массив длинною в 2 элемента
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if headerParts[1] == "" {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1]) // парсим jwt токен и обрабатываем ошибку
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId) // Записываем значение id пользователя в context для того что бы иметь доступ к id пользователя в последующих обработчиках которые вызываются после этой прослойки
}

func getUserId(c *gin.Context) (int, error) { // Функция для получения id пользователя
	id, ok := c.Get(userCtx) // Возвращает Id пользователя типа interface
	if !ok {                 // Проверяем существует ли id
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int) // Приводим id к типу int
	if !ok {              // Если приведение к типу int не удалось возвращаем ошибку
		return 0, errors.New("user id is invalid type")
	}

	return idInt, nil
}
