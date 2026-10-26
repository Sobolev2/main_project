package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	kafkago "github.com/segmentio/kafka-go"
)

type MessageSent struct {
	Event      string `json:"event"`
	MessageID  int    `json:"message_id"`
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
}

type NotificationHandler interface {
	CreateMessageNotification(
		ctx context.Context,
		userID, senderID, messageID int,
	) error
}

func StartConsumer(
	ctx context.Context,
	broker string,
	handler NotificationHandler,
) error {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "message-events",
		GroupID: "notification-service",
	})

	defer func() {
		if err := reader.Close(); err != nil {
			slog.Error("failed to close Kafka reader", "error", err)
		}
	}()

	slog.Info("Kafka consumer started", "topic", "message-events")

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("fetch Kafka message: %w", err)
		}

		var event MessageSent

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("decode Kafka event: %w", err)
		}

		if event.Event != "message.sent" {
			return fmt.Errorf("unexpected event type: %s", event.Event)
		}

		// Создаём уведомление через Controller.
		if err := handler.CreateMessageNotification(
			ctx,
			event.ReceiverID,
			event.SenderID,
			event.MessageID,
		); err != nil {
			return fmt.Errorf("create notification: %w", err)
		}

		// Подтверждаем событие только после сохранения.
		if err := reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("commit Kafka message: %w", err)
		}

		slog.Info(
			"Kafka event processed successfully",
			"message_id", event.MessageID,
			"sender_id", event.SenderID,
			"receiver_id", event.ReceiverID,
		)
	}
}
