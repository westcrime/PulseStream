package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"

	"github.com/westcrime77/PulseStream/pkg/config"
)

const MaxBytes = 10e6

type Reader struct {
	reader *kafka.Reader
}

func NewReader(lc fx.Lifecycle, cfg config.Kafka) *Reader {
	r := &Reader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
			GroupID:  cfg.GroupID,
			Topic:    cfg.Topic,
			MaxBytes: MaxBytes,
		}),
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return r.reader.Close()
		},
	})

	return r
}

func (r *Reader) Fetch(ctx context.Context) (kafka.Message, error) {
	return r.reader.FetchMessage(ctx)
}

func (r *Reader) Commit(ctx context.Context, msg kafka.Message) error {
	return r.reader.CommitMessages(ctx, msg)
}

func (r *Reader) Close() error {
	return r.reader.Close()
}
