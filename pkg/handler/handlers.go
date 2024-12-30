package handler

import (
	"github.com/gin-gonic/gin" // используется gin web framework
	"github.com/speeddem0n/todoapp/pkg/service"
)

type Handler struct { // Структура handler
	services *service.Service
}

func NewHandler(services *service.Service) *Handler { // Инициализируем новую структуру handler и передаем в нее стуруктуру Service с уровня .\services (2 УРОВЕНЬ)
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // Устанавлевает ReleaseMod Для запуска сервера
	router := gin.New()          // Инициализация роутера

	auth := router.Group("/auth") // Группа авторизации
	{
		auth.POST("/sign-up", h.signUp) // Регистрация
		auth.POST("/sign-in", h.signIn) // Авторизация
	}

	api := router.Group("/api", h.userIdentity) // Группа API
	{
		lists := api.Group("/lists") // Группа списков
		{
			lists.POST("/", h.createList)      // Создать список
			lists.GET("/", h.getAllLists)      // Получить все списки
			lists.GET("/:id", h.getListById)   // Получить список по Id
			lists.PUT("/:id", h.updateList)    // Обновить список по Id
			lists.DELETE("/:id", h.deleteList) // Удалить список по Id

			items := lists.Group(":id/items") // Группа элементов списка
			{
				items.POST("/", h.createItem) // Создать элемент списка
				items.GET("/", h.getAllItems) // Получить все элементы списка
			}
		}

		items := api.Group("items") // Группа элементов списка
		{
			items.GET("/:id", h.getItemById)   // Получить элемент списка по Id
			items.PUT("/:id", h.updateItem)    // Обновить элемент списка по Id
			items.DELETE("/:id", h.deleteItem) // Удалить элемент списка по Id
		}
	}

	return router
}
