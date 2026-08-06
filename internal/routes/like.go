package routes

import (
	controller "semen_project/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupLikeRoutes(r *gin.RouterGroup, h *controller.LikeHandler) {
	r.POST("/posts/:id/like", h.LikePost)
	r.DELETE("/posts/:id/like", h.DeleteLike)
	r.GET("/posts/:id/likes/count", h.GetCountLikes)
	r.GET("/posts/:id/likes", h.GetAllPostLikes)
	r.GET("/users/:id/likes", h.GetAllUserLikes)
	r.GET("/posts/:id/like-status", h.GetLikeStatus)
}