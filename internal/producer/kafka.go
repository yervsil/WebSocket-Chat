package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	w *kafka.Writer
}


func NewProducer(brokers []string, topic string) *KafkaProducer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &KafkaProducer{w: w}
}

func (p *KafkaProducer) ProduceJSONMessage(ctx context.Context, data any) error {

	payload, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("can't marshal data: %w", err)
	}

	err = p.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(uuid.New().String()),
		Value: payload,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.w.Close()
}