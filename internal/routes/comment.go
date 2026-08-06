package routes

import (
	controller "semen_project/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupCommentRoutes(r *gin.RouterGroup, h *controller.CommentHandler) {
	r.POST("/posts/:id/comments", h.CreateComment)
	r.DELETE("/comments/:id", h.DeleteComment)
	r.PATCH("/comments/:id", h.UpdateComment)
	r.GET("/posts/:id/comments", h.GetAllPostComments)
	r.GET("/posts/:id/comments/count", h.GetCountComments)
}