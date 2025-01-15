package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin" // используется gin web framework
	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/internal/models"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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

	id, err := h.services.TodoItem.Create(userId, listId, input) // Вызываем метод Create для создания нового элемента списка
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{ // В ответ возвращаем id только что созданного "todo" (gin.H тоже самое что map[string]interface{})
		"id": id,
	})
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	itemId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateItemInput
	err = c.BindJSON(&input) // Получаем инпут от пользователя и записываем его в структуру input todo.UpdateListInput
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.TodoItem.Update(userId, itemId, input) // Вызывает метод Delete из сервисов для обновления списка по listID
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	}) // Возващаем Структуру statusResponse и пишем в ней что status: ok
}

func (h *Handler) deleteItem(c *gin.Context) { // Метод обработчика для удаления эелемнта списка по его ID
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	itemId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	err = h.services.TodoItem.Delete(userId, itemId) // Вызывает метод Delete из сервисов для удаления элемента списка по его id
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	}) // Возващаем Структуру statusResponse и пишем в ней что status: ok

}
