package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/westcrime77/PulseStream/pkg/config"
)

const MaxBytes = 10e6

type Reader struct {
	reader *kafka.Reader
}

func NewReader(cfg config.Kafka) *Reader {
	return &Reader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
			GroupID:  cfg.GroupID,
			Topic:    cfg.Topic,
			MaxBytes: MaxBytes,
		}),
	}
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
