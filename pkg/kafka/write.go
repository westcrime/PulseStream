package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/westcrime77/PulseStream/pkg/config"
	"go.uber.org/fx"
)

type Writer struct {
	writer *kafka.Writer
}

func NewWriter(lc fx.Lifecycle, cfg config.Kafka) *Writer {
	w := &Writer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)),
			Topic:                  cfg.Topic,
			AllowAutoTopicCreation: true,
			BatchTimeout:           500 * time.Millisecond,
		},
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return w.writer.Close()
		},
	})

	return w
}

func (w *Writer) Write(ctx context.Context, messages []kafka.Message) error {
	return w.writer.WriteMessages(ctx, messages...)
}

func (w *Writer) Close() error {
	return w.writer.Close()
}
