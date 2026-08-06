package routes

import (
	controller "semen_project/internal/controllers"
	"semen_project/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *controller.Handlers, secret string) {

    api := r.Group("/api")

    public := api.Group("")
    private := api.Group("")

    private.Use(middleware.AuthMiddleware(secret))

    SetupPublicAuthenticationRoutes(public, h.AuthHandler)

    SetupPrivateAuthenticationRoutes(private, h.AuthHandler)

    SetupUserRoutes(private, h.UserHandler)
    SetupPostRoutes(private, h.PostHandler)
    SetupMessageRoutes(private, h.MessageHandler)
    SetupLikeRoutes(private, h.LikeHandler)
    SetupFriendRoutes(private, h.FriendHandler)
    SetupFollowRoutes(private, h.FollowHandler)
    SetupCommentRoutes(private, h.CommentHandler)
    SetupChatRoutes(private, h.ChatHandler)
}