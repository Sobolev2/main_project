package routes

import (
	controller "semen_project/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupFollowRoutes(r *gin.RouterGroup, h *controller.FollowHandler) {
	r.POST("/:id/follow", h.FollowUser)
	r.DELETE("/:id/follow", h.UnFollowUser)
	r.GET("/:id/followers", h.GetAllFollowers)
	r.GET("/:id/following", h.GetAllFollowing)
	r.GET("/:id/followers/count", h.GetCountFollowers)
	r.GET("/:id/following/count", h.GetCountFollowing)
	r.GET("/:id/follow-status", h.CheckFollowStatus)
}