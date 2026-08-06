package routes

import (
	controller "semen_project/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupPostRoutes(r *gin.RouterGroup, h *controller.PostHandler) {
	r.POST("/posts", h.CreatePost)
	r.GET("/posts/:id", h.GetPostById)
	r.GET("/feed", h.GetFeed)
	r.GET("/users/:id/posts", h.GetAllUserPosts)
	r.GET("/posts", h.GetAllPosts)
	r.DELETE("/posts/:id", h.DeletePost)
	r.PATCH("/posts/:id", h.UpdatePost)
}