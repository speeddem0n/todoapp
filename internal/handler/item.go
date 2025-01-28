package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/speeddem0n/todoapp/internal/models"
)

// @Summary Create todo list item
// @Security ApiKeyAuth
// @Tags items
// @Description Create new todo list item
// @Accept  json
// @Produce  json
// @Param id path int true "list ID"
// @Param input body models.CreateTodoItemInp true "item info"
// @Success 201 {integer} integer itemID
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id}/items [post]
func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id")) // достаем id списка из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error on getting list id: %v", err))
		return
	}

	var input models.CreateTodoItemInp
	err = c.BindJSON(&input) // Считываем инпут пользователя
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalit input: %v", err))
		return
	}

	id, err := h.services.TodoItem.Create(userId, listId, input) // Вызываем метод Create для создания нового элемента списка
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error on creating list item: %v", err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{ // В ответ возвращаем id только что созданного "todo" (gin.H тоже самое что map[string]interface{})
		"id": id,
	})
}

// @Summary Get All Lists items
// @Security ApiKeyAuth
// @Tags items
// @Description Get All todo list items
// @Accept  json
// @Produce  json
// @Param id path int true "list ID"
// @Success 200 {array} models.TodoItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id}/items [get]
func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error on getting list id: %v", err))
		return
	}

	items, err := h.services.TodoItem.GetAll(userId, listId) // Вызывает метод GetAll из сервисов для получения всех элементов списка
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error on getting items: %v", err))
		return
	}

	c.JSON(http.StatusOK, items) // Передаем слайс с элементами списка в тело ответа

}

// @Summary Get List Item By Id
// @Security ApiKeyAuth
// @Tags items
// @Description Get list item using itemID
// @Accept  json
// @Produce  json
// @Param id path int true "item ID"
// @Success 200 {object} models.TodoItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/{id} [get]
func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	itemId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error on getting item id: %v", err))
		return
	}

	item, err := h.services.TodoItem.GetById(userId, itemId) // Вызывает метод GetById из сервисов для получения элемента списка по его id
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error on getting item by id: %d, %v", itemId, err))
		return
	}

	c.JSON(http.StatusOK, item) // Записываем нужный элемент списка в тело ответа
}

// @Summary Update List item By Id
// @Security ApiKeyAuth
// @Tags items
// @Description Update list item using itemID
// @Accept  json
// @Produce  json
// @Param input body models.UpdateItemInput true "Update item input"
// @Param id path int true "item ID"
// @Success 200 {string} string "status"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/{id} [put]
func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	itemId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error on getting item id: %v", err))
		return
	}

	var input models.UpdateItemInput
	err = c.BindJSON(&input) // Получаем инпут от пользователя и записываем его в структуру input todo.UpdateListInput
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalit input: %v", err))
		return
	}

	err = h.services.TodoItem.Update(userId, itemId, input) // Вызывает метод Delete из сервисов для обновления списка по listID
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error on updating item by id: %d, %v", itemId, err))
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	}) // Возващаем Структуру statusResponse и пишем в ней что status: ok
}

// @Summary Delete List Item By Id
// @Security ApiKeyAuth
// @Tags items
// @Description Delete list item using itemID
// @Accept  json
// @Produce  json
// @Param id path int true "item ID"
// @Success 200 {string} string "status"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/{id} [delete]
func (h *Handler) deleteItem(c *gin.Context) { // Метод обработчика для удаления эелемнта списка по его ID
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	itemId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalit input: %v", err))
		return
	}

	err = h.services.TodoItem.Delete(userId, itemId) // Вызывает метод Delete из сервисов для удаления элемента списка по его id
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("invalit input: %v", err))
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	}) // Возващаем Структуру statusResponse и пишем в ней что status: ok

}
