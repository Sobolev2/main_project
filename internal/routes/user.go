package routes

import (
	controller "semen_project/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup, h *controller.UserHandler) {
	r.GET("/me", h.GetMe)
	r.PUT("/user", h.UpdateUser)
	r.GET("/users", h.GetAllUsers)
	r.DELETE("/user", h.DeleteUser)
	r.POST("/password", h.UpdatePassword)
	r.GET("/user/:id", h.GetUserById)
	r.GET("/users/search", h.SearchUsers)
	r.GET("/user/:username", h.GetUserByUsername)
}
