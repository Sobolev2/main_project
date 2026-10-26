package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "semen_project/proto/gen/notificationpb"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	client pb.NotificationServiceClient
}

func NewNotificationHandler(
	client pb.NotificationServiceClient,
) *NotificationHandler {
	return &NotificationHandler{
		client: client,
	}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetInt("userID")

	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	response, err := h.client.GetNotifications(
		ctx,
		&pb.GetNotificationsRequest{
			UserId: int64(userID),
			Limit:  20,
			Offset: 0,
		},
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "notification service unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetInt("userID")
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	response, err := h.client.GetUnreadCount(ctx, &pb.GetUnreadCountRequest{
		UserId: int64(userID),
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "notification service unavailable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": response.GetCount()})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetInt("userID")
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	notificationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || notificationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	response, err := h.client.MarkAsRead(ctx, &pb.MarkAsReadRequest{
		UserId:         int64(userID),
		NotificationId: notificationID,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "notification service unavailable"})
		return
	}

	if !response.GetSuccess() {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
