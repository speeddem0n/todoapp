package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin" // используется gin вместо стандартого net/http
	todo "github.com/speeddem0n/todoapp"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id")) // достаем id списка из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	var input todo.TodoItem
	err = c.BindJSON(&input) // Считываем инпут пользователя в input
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoItem.Create(userId, listId, input) // вызываем метод Create для создания нового элемента списка
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{ // В ответ возвращаем id только что созданного "todo"
		"id": id,
	})
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	items, err := h.services.TodoItem.GetAll(userId, listId) // Вызывает метод GetAll из сервисов для получения всех элементов списка
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items) // Передаем слайс с элементами списка в тело ответа

}

func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	item, err := h.services.TodoItem.GetById(userId, itemId) // Вызывает метод GetById из сервисов для получения элемента списка по его id
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item) // Записываем нужный элемент списка в тело ответа
}

func (h *Handler) updateItem(c *gin.Context) {

}

func (h *Handler) deleteItem(c *gin.Context) {

}
