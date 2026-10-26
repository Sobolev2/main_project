package routes

import (
	"semen_project/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupNotificationRoutes(
	r *gin.RouterGroup,
	h *controllers.NotificationHandler,
) {
	r.GET("/notifications", h.GetNotifications)
	r.GET("/notifications/unread/count", h.GetUnreadCount)
	r.PATCH("/notifications/:id/read", h.MarkAsRead)
}
