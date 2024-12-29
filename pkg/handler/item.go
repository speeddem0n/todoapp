package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin" // используется gin вместо стандартого net/http
	todo "github.com/speeddem0n/todoapp"
)

// @Summary		Create todo item
// @Security		ApiKeyAuth
// @Tags			items
// @Description	create todo item
// @ID				create-item
// @Accept			json
// @Produce		json
// @Param			input	body		todo.TodoItem	true	"item info"
// @Success		200		{integer}	integer			1
// @Failure		400,404	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/lists/items [post]
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

	id, err := h.services.TodoItem.Create(userId, listId, input) // вызываем метод Create для создания нового элемента списка
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{ // В ответ возвращаем id только что созданного "todo" (gin.H тоже самое что map[string]interface{})
		"id": id,
	})
}

// @Summary		Get All Todo Items
// @Security		ApiKeyAuth
// @Tags			items
// @Description	get all todo items
// @ID				get-all-items
// @Accept			json
// @Produce		json
// @Success		200		{object}	[]todo.TodoItem	TodoItems
// @Failure		400,404	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/lists/items [get]
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

// @Summary		Get Todo Item By Id
// @Security		ApiKeyAuth
// @Tags			items
// @Description	get todo item by id
// @ID				get-item-by-id
// @Accept			json
// @Produce		json
// @Success		200		{object}	todo.TodoItem
// @Failure		400,404	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/items/:id [get]
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

// @Summary		Update Todo Item
// @Security		ApiKeyAuth
// @Tags			items
// @Description	update todo item
// @ID				update-item
// @Accept			json
// @Produce		json
// @Param			input	body		todo.UpdateItemInput	true	"update item info"
// @Success		200		{string}	string					"ok"
// @Failure		400,404	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/items/:id [put]
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

	var input todo.UpdateItemInput
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
	}) // Возващаем Структуру statusResponse и пишем в ней что все ok
}

// @Summary		Delete Todo Item
// @Security		ApiKeyAuth
// @Tags			items
// @Description	delete todo item
// @ID				delete-item
// @Accept			json
// @Produce		json
// @Success		200		{string}	string	"ok"
// @Failure		400,404	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/items/:id [delete]
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
	}) // Возващаем Структуру statusResponse и пишем в ней что все ok

}
