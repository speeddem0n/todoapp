package handler

import (
	"github.com/gin-gonic/gin" // используется gin вместо стандартого net/http
	"github.com/speeddem0n/todoapp/pkg/service"
)

type Handler struct { // Структура handler
	services *service.Service
}

func NewHandler(services *service.Service) *Handler { // Инициализируем новую структуру handler и передаем в нее стуруктуру Service с уровня .\services (2 УРОВЕНЬ)
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New() // Инициализация роутера

	auth := router.Group("/auth") // Группа авторизации
	{
		auth.POST("/sign-up", h.singUp) // Регистрация
		auth.POST("/sign-in", h.singIn) // Авторизация
	}

	api := router.Group("/api", h.userIdentity) // Группа API
	{
		lists := api.Group("/lists") // Группа списков
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
				items.GET("/:item_id", h.getItemById)
				items.PUT("/:item_id", h.updateItem)
				items.DELETE("/:item_id", h.deleteItem)
			}
		}
	}

	return router
}
