package handler

import (
	"net/http"

	"github.com/gin-gonic/gin" // используется gin вместо стандартого net/http
	todo "github.com/speeddem0n/todoapp"
)

func (h *Handler) singUp(c *gin.Context) { // Метод обработчик для Регистрации
	var input todo.User

	err := c.BindJSON(&input) // BindJSON принимает ссылку на объект в который мы хотим распарсить тело JSON
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error()) // Возвращается ошибка 400 (Не корректные данные в запросе от пользователя)
		return
	}

	id, err := h.services.Authorization.CreateUser(input) // Вызываем из сервисов метод для создания пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) // В случае ошибки возвращаем код InternalServerError
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) singIn(c *gin.Context) { // Метод обработчик для Авторизации

}
