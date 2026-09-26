package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	kafkago "github.com/segmentio/kafka-go"
	"semen_project/internal/events"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer() *Producer {
	return &Producer{
		writer: &kafkago.Writer{
			Addr:     kafkago.TCP("localhost:9092"),
			Topic:    "post-events",
			Balancer: &kafkago.Hash{},
		},
	}
}

func (p *Producer) PublishPostLiked(ctx context.Context, event events.PostLiked) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(fmt.Sprintf("post-%d", event.PostID)),
		Value: value,
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
