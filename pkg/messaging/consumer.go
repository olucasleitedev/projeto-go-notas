package messaging

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, key, value []byte) error

func Consume(ctx context.Context, brokers []string, group, topic string, handler MessageHandler) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: group,
		Topic:   topic,
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := handler(ctx, msg.Key, msg.Value); err != nil {
			return err
		}
	}
}

func UnmarshalJSON(value []byte, dest any) error {
	return json.Unmarshal(value, dest)
}
