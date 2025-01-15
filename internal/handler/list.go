package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin" // используется gin вместо стандартого net/http
	"github.com/speeddem0n/todoapp/internal/models"
)

// @Summary Create todo list
// @Security ApiKeyAuth
// @Tags lists
// @Description Create new todo list
// @Accept  json
// @Produce  json
// @Param input body models.TodoList true "list info"
// @Success 200 {integer} integer listID
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input models.TodoList
	err = c.BindJSON(&input) // Считываем инпут пользователя в input
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoList.Create(userId, input) // Вызываем метод Create для создания нового списка
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{ // В ответ возвращаем id только что созданного списка (gin.H тоже самое что map[string]interface{})
		"id": id,
	})
}

// @Summary Get All Lists
// @Security ApiKeyAuth
// @Tags lists
// @Description Get all user lists
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllListsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	lists, err := h.services.TodoList.GetAll(userId) // Вызывает метод GetALL из сервисов для получения всех списков пользоваетля
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListsResponse{ // Записываем в тело ответа структуру getAllListsResponse которая содержит слайс Списков пользователя
		Data: lists,
	})

}

// @Summary Get List By Id
// @Security ApiKeyAuth
// @Tags lists
// @Description Get list using listID
// @Accept  json
// @Produce  json
// @Param id path int true "list ID"
// @Success 200 {object} models.ListItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id} [get]
func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.TodoList.GetById(userId, listId) // Вызывает метод GetById из сервисов для получения списка по его id
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list) // Записываем нужный список в тело ответа

}

// @Summary Update List By Id
// @Security ApiKeyAuth
// @Tags lists
// @Description Update list using listID
// @Accept  json
// @Produce  json
// @Param id path int true "list ID"
// @Success 200 {string} string "status"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id} [put]
func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateListInput
	err = c.BindJSON(&input) // Получаем инпут от пользователя и записываем его в структуру input todo.UpdateListInput
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.TodoList.Update(userId, listId, input) // Вызывает метод Delete из сервисов для обновления списка по listID
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	}) // Возващаем Структуру statusResponse и пишем в ней status: ok
}

// @Summary Delete List By Id
// @Security ApiKeyAuth
// @Tags lists
// @Description Delete list using listID
// @Accept  json
// @Produce  json
// @Param id path int true "list ID"
// @Success 200 {string} string "status"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return

	}

	err = h.services.TodoList.Delete(userId, listId) // Вызывает метод Delete из сервисов для удаления списка по listID
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	}) // Возващаем Структуру statusResponse и пишем в ней status: ok

}
