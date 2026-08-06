package routes

import (
	controller "semen_project/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupMessageRoutes(r *gin.RouterGroup, h *controller.MessageHandler) {
	r.POST("/chats/:id/messages", h.SendMessage)
	r.PATCH("/messages/:id", h.UpdateMessage)
	r.DELETE("/messages/:id", h.DeleteMessage)
	r.GET("/chats/:id/messages", h.GetAllChatMessages)
	r.PATCH("/messages/:id/status/read", h.UpdateMarkReadToRead)
	r.GET("/chats/:id/messages/unread/count", h.GetCountNotReadMessages)
}