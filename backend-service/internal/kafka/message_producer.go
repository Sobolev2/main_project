package kafka

import (
	"context"
	"encoding/json"
	"strconv"

	kafkago "github.com/segmentio/kafka-go"
	"semen_project/internal/events"
)

type MessageProducer struct {
	writer *kafkago.Writer
}

func NewMessageProducer(broker string) *MessageProducer {
	return &MessageProducer{
		writer: &kafkago.Writer{
			Addr:     kafkago.TCP(broker),
			Topic:    "message-events",
			Balancer: &kafkago.Hash{},
		},
	}
}

func (p *MessageProducer) PublishMessageSent(
	ctx context.Context,
	event events.MessageSent,
) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(strconv.Itoa(event.ReceiverID)),
		Value: data,
	})
}

func (p *MessageProducer) Close() error {
	return p.writer.Close()
}
