package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) { // Функция для обработки http ошибок
	logrus.Error(message)                                              // Выводим сообщение об ошибке в консоль
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message}) // AbortWithStatusJSON принимает hhtp status code и тело ответа. AbortWithStatusJSON блокирует выполнение последующих обработчиков и записывает в ответ статус код и тело сообщения в формате JSON
}
