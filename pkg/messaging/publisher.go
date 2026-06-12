package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"estudos-golang/pkg/events"

	"github.com/segmentio/kafka-go"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, key string, payload any) error
	Close() error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(context.Context, string, string, any) error { return nil }
func (NoopPublisher) Close() error                                       { return nil }

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, topic string, key string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: body,
	})
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}

func PublishNoteEvent(ctx context.Context, pub Publisher, evt events.NoteEvent) error {
	if pub == nil {
		return nil
	}
	return pub.Publish(ctx, events.TopicNoteEvents, evt.NoteID, evt)
}
