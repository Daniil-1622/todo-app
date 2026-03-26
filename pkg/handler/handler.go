package handler

import (
	"github.com/Daniil-1622/todo-app/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Здесь описывается структура маршрутов (роутинг) для todo-приложения.
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp) // Регистрация
		auth.POST("/sign-in", h.signIn) // Вход
	}
	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{ // Маршруты API /api:
			lists.POST("/", h.createList)      // Создать список
			lists.GET("/", h.getAllList)       // Получить все списки
			lists.GET("/:id", h.getListById)   // Получить список по id
			lists.PUT("/:id", h.updateList)    // Обновить список
			lists.DELETE("/:id", h.deleteList) // Удалить список

			items := lists.Group((":id/items"))
			{ // Элементы списка /api/lists/:id/items:
				items.POST("/", h.createItem)            // Создать элемент списка
				items.GET("/", h.getAllItem)             // Получить все элементы списка
				items.GET("/:items_id", h.getItemById)   // Получить элемент списка по id
				items.PUT("/:items_id", h.updateItem)    // Обновить элемент списка
				items.DELETE("/:items_id", h.deleteItem) // Удалить элемент списка
			}
		}
	}
	return router
}
