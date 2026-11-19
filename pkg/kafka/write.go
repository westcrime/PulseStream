package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/westcrime77/PulseStream/pkg/config"
)

type Writer struct {
	writer *kafka.Writer
}

func NewWriter(cfg config.Kafka) *Writer {
	return &Writer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)),
			Topic:                  cfg.Topic,
			AllowAutoTopicCreation: true,
		},
	}
}

func (w *Writer) Write(ctx context.Context, messages []kafka.Message) error {
	return w.writer.WriteMessages(ctx, messages...)
}

func (w *Writer) Close() error {
	return w.writer.Close()
}
