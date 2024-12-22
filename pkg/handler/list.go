package handler

import (
	"net/http"

	"github.com/gin-gonic/gin" // используется gin вместо стандартого net/http
	todo "github.com/speeddem0n/todoapp"
)

func (h *Handler) createList(c *gin.Context) { // Метод обработчика для создания Списка todo
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		return
	}

	var input todo.TodoList
	err = c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllListsResponse struct {
	Data []todo.TodoList `json:"data"`
}

func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetALL(userId) // Вызывает метод GetALL из сервисов для получения всех списков пользоваетля
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListsResponse{ // Записываем в тело ответа структуру getAllListsResponse которая содержит слайс Списков пользователя
		Data: lists,
	})

}

func (h *Handler) getListById(c *gin.Context) {

}

func (h *Handler) updateList(c *gin.Context) {

}

func (h *Handler) deleteList(c *gin.Context) {

}
