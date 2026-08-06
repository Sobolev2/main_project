package routes

import (
	controller "semen_project/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupFriendRoutes(r *gin.RouterGroup, h *controller.FriendHandler) {
	r.GET("/:id/friends", h.GetAllFriends)
	r.GET("/friends", h.GetAllMyFriends)
	r.GET("/:id/friends/count", h.GetCountFriends)
	r.GET("/:id/friend-status", h.CheckFriendship)
}