package routes

import (
	controller "semen_project/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(r *gin.RouterGroup, h *controller.ChatHandler) {
	r.POST("/users/:id/chat", h.CreateChat)
	r.GET("/chats", h.GetAllUserChats)
	r.DELETE("/chats/:id", h.DeleteChat)
	r.GET("/users/:id/chat", h.GetChatByUserID)
}