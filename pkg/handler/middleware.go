package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) { // метод идентификации пользователя по jwt токену
	header := c.GetHeader(authorizationHeader) // Поулчаем header
	if header == "" {                          // Проверяем что хедер не пустой
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ") // Разделяем строку header по пробелам
	if len(headerParts) != 2 {                // При корректном формате хедера strings.Split должен вернуть массив длинною в 2 элемента
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1]) // парсим jwt токен и обрабатываем ошибку
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId) // Записываем значение id пользователя в context для того что бы иметь доступ к id пользователя в последующих обработчиках которые вызываются после этой прослойки
}
