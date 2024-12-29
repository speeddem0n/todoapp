package handler

import (
	"net/http"

	"github.com/gin-gonic/gin" // используется gin web framework
	todo "github.com/speeddem0n/todoapp"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body todo.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) singUp(c *gin.Context) { // Метод обработчик для Регистрации пользователей
	var input todo.User

	err := c.BindJSON(&input) // BindJSON принимает ссылку на объект в который мы хотим распарсить тело JSON
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body") // Возвращается ошибка 400 (Не корректные данные в запросе от пользователя)
		return
	}

	id, err := h.services.Authorization.CreateUser(input) // Вызываем из сервисов метод для создания пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) // В случае ошибки возвращаем код InternalServerError
		return
	}

	c.JSON(http.StatusOK, gin.H{ // Передаем id бользоваетля в случае успеха (gin.H тоже самое что map[string]interface{})
		"id": id,
	})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept json
// @Produce json
// @Param input body signInInput true "username & password"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
type signInInput struct { // Структура для sign in юзера
	Username string `json:"username" binding:"required"` // Имя пользователя
	Password string `json:"password" binding:"required"` // Пароль
}

func (h *Handler) singIn(c *gin.Context) { // Метод обработчик для Авторизации
	var input signInInput // Объявляем пустую структуру signInInput

	err := c.BindJSON(&input) // BindJSON принимает ссылку на объект в который мы хотим распарсить тело JSON
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body") // Возвращается ошибка 400 (Не корректные данные в запросе от пользователя)
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password) // Вызываем из сервисов метод для получения JWT токена пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) // В случае ошибки возвращаем код InternalServerError
		return
	}

	c.JSON(http.StatusOK, gin.H{ // Передаем token пользоваетля в случае успеха (gin.H тоже самое что map[string]interface{})
		"token": token,
	})
}
