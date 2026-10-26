package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Notification struct {
	ID        int64
	UserID    int64
	SenderID  int64
	MessageID int64
	Type      string
	IsRead    bool
	CreatedAt time.Time
}

type NotificationRepo interface {
	CreateNotification(
		ctx context.Context,
		userID, senderID, messageID int,
	) error

	GetNotifications(
		ctx context.Context,
		userID int64,
		limit, offset int32,
	) ([]Notification, error)
	GetUnreadCount(ctx context.Context, userID int64) (int64, error)
	MarkAsRead(ctx context.Context, userID, notificationID int64) (bool, error)
}

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) CreateNotification(
	ctx context.Context,
	userID, senderID, messageID int,
) error {
	query := `
		INSERT INTO notifications (
			user_id,
			sender_id,
			message_id,
			type
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (type, message_id) DO NOTHING
	`

	_, err := s.db.Exec(
		ctx,
		query,
		userID,
		senderID,
		messageID,
		"new_message",
	)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}

	return nil
}

func (s *Store) GetNotifications(
	ctx context.Context,
	userID int64,
	limit, offset int32,
) ([]Notification, error) {
	query := `
		SELECT
			id,
			user_id,
			sender_id,
			message_id,
			type,
			is_read,
			created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(
		ctx,
		query,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get notifications: %w", err)
	}
	defer rows.Close()

	notifications := make([]Notification, 0)

	for rows.Next() {
		var n Notification

		if err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.SenderID,
			&n.MessageID,
			&n.Type,
			&n.IsRead,
			&n.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}

		notifications = append(notifications, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notifications: %w", err)
	}

	return notifications, nil
}
func (s *Store) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	var count int64

	err := s.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = false
	`, userID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("get unread count: %w", err)
	}

	return count, nil
}

func (s *Store) MarkAsRead(ctx context.Context, userID, notificationID int64) (bool, error) {
	result, err := s.db.Exec(ctx, `
		UPDATE notifications
		SET is_read = true
		WHERE id = $1 AND user_id = $2
	`, notificationID, userID)

	if err != nil {
		return false, fmt.Errorf("mark notification as read: %w", err)
	}

	return result.RowsAffected() > 0, nil
}
