package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin" // используется gin web framework
	"github.com/speeddem0n/todoapp/internal/models"
)

// @Summary SignUp
// @Tags auth
// @Description Create new account
// @Accept  json
// @Produce  json
// @Param input body models.User true "Account info"
// @Success 200 {integer} integer userId
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) { // Метод обработчик для Регистрации пользователей
	var input models.User

	err := c.BindJSON(&input) // BindJSON принимает ссылку на объект в который мы хотим распарсить тело JSON
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input") // Возвращается ошибка 400 (Не корректные данные в запросе от пользователя)
		return
	}

	id, err := h.services.Authorization.CreateUser(input) // Вызываем из сервисов метод для создания пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("failed to sign in: %v", err)) // В случае ошибки возвращаем код InternalServerError
		return
	}

	c.JSON(http.StatusOK, gin.H{ // Передаем id пользоваетеля в случае успеха (gin.H тоже самое что map[string]interface{})
		"id": id,
	})
}

type signInInput struct { // Структура для sign in юзера
	Username string `json:"username" binding:"required"` // Имя пользователя
	Password string `json:"password" binding:"required"` // Пароль
}

// @Summary SignIn
// @Tags auth
// @Description login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "username and passwor"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) { // Метод обработчик для Авторизации
	var input signInInput // Объявляем пустую структуру signInInput

	err := c.BindJSON(&input) // BindJSON принимает ссылку на объект в который мы хотим распарсить тело JSON
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input") // Возвращается ошибка 400 (Некорректные данные в запросе от пользователя)
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password) // Вызываем из сервисов метод для получения JWT токена пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("failed to sign ip: %v", err)) // В случае ошибки возвращаем код InternalServerError
		return
	}

	c.JSON(http.StatusOK, gin.H{ // Передаем token пользоваетеля в случае успеха (gin.H тоже самое что map[string]interface{})
		"token": token,
	})
}
