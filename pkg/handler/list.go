package handler

import (
	"net/http"
	"strconv"

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

type getAllListsResponse struct { // Структура для записи слайса списков что бы потом передать ее в тело ответа метода getAllLists()
	Data []todo.TodoList `json:"data"`
}

func (h *Handler) getAllLists(c *gin.Context) { // Метод для получения всех списков конкретного пользователя
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

func (h *Handler) getListById(c *gin.Context) { // Метод для получения конкретного списка пользователя по его ID
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id")) // достаем id из URL param
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return

	}

	list, err := h.services.TodoList.GetById(userId, listId) // Вызывает метод GetALL из сервисов для получения всех списков пользоваетля
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list) // Записываем нужный список в тело ответа

}

func (h *Handler) updateList(c *gin.Context) {

}

func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c) // Обращаемся к функции getUserId из middleware для получения id пользователя
	if err != nil {
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
	}) // Возващаем Структуру statusResponse и пишем в ней что все ok

}
