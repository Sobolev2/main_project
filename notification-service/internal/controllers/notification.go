package controllers

import (
	"context"
	"fmt"
	"log/slog"

	"notification-service/internal/repository"
)

type NotificationHandler struct {
	Repo repository.NotificationRepo
}

func NewNotificationHandler(
	repo repository.NotificationRepo,
) *NotificationHandler {
	return &NotificationHandler{
		Repo: repo,
	}
}

func (h *NotificationHandler) CreateMessageNotification(
	ctx context.Context,
	userID, senderID, messageID int,
) error {
	// Проверяем корректность данных.
	if userID <= 0 || senderID <= 0 || messageID <= 0 {
		return fmt.Errorf("invalid notification data")
	}

	if userID == senderID {
		return fmt.Errorf("sender and receiver cannot be the same")
	}

	// Сохраняем уведомление в PostgreSQL.
	if err := h.Repo.CreateNotification(
		ctx,
		userID,
		senderID,
		messageID,
	); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	slog.Info(
		"Notification saved successfully",
		"user_id", userID,
		"sender_id", senderID,
		"message_id", messageID,
	)

	return nil
}

func (h *NotificationHandler) GetNotifications(
	ctx context.Context,
	userID int64,
	limit, offset int32,
) ([]repository.Notification, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	notifications, err := h.Repo.GetNotifications(
		ctx,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return notifications, nil
}
func (h *NotificationHandler) GetUnreadCount(
	ctx context.Context,
	userID int64,
) (int64, error) {
	if userID <= 0 {
		return 0, fmt.Errorf("invalid user ID")
	}

	return h.Repo.GetUnreadCount(ctx, userID)
}

func (h *NotificationHandler) MarkAsRead(
	ctx context.Context,
	userID, notificationID int64,
) (bool, error) {
	if userID <= 0 || notificationID <= 0 {
		return false, fmt.Errorf("invalid notification data")
	}

	return h.Repo.MarkAsRead(ctx, userID, notificationID)
}
