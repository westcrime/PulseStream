package batcher

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/westcrime77/PulseStream/pkg/kafka"
	kafkamodels "github.com/westcrime77/PulseStream/pkg/kafka/models"
	"github.com/westcrime77/PulseStream/services/ingestion/internal/http/models"
	"github.com/westcrime77/PulseStream/services/ingestion/internal/metrics"
	"go.uber.org/fx"
)

const BufferSize = 2048
const CheckBufferInterval = 1 * time.Second
const WriteTimeout = 5 * time.Second

type Logger interface {
	Errorf(format string, v ...any)
}

type Batcher struct {
	events   chan models.Event
	writer   *kafka.Writer
	log      Logger
	stop     chan struct{}
	done     chan struct{}
	stopOnce sync.Once
}

func (b *Batcher) Run() {
	defer close(b.done)
	ticker := time.NewTicker(CheckBufferInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			b.processBuffer()
		case <-b.stop:
			b.processBuffer()
			return
		}
	}
}

func (b *Batcher) Stop(ctx context.Context) error {
	b.stopOnce.Do(func() { close(b.stop) })
	select {
	case <-b.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func New(lc fx.Lifecycle, writer *kafka.Writer, log Logger) *Batcher {
	b := &Batcher{
		events: make(chan models.Event, BufferSize),
		writer: writer,
		log:    log,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go b.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return b.Stop(ctx)
		},
	})

	return b
}

func (b *Batcher) Add(events ...models.Event) (rejected []models.Event) {
	for i, e := range events {
		select {
		case b.events <- e:
		default:
			return events[i:]
		}
	}
	return nil
}

func (b *Batcher) drain() []models.Event {
	batch := make([]models.Event, 0, BufferSize)
	for {
		select {
		case ev := <-b.events:
			batch = append(batch, ev)
		default:
			return batch
		}
	}
}

func (b *Batcher) processBuffer() {
	batch := b.drain()
	if len(batch) == 0 {
		return
	}
	b.PublishEvents(batch)
}

func (b *Batcher) PublishEvents(events []models.Event) {
	ctx, cancel := context.WithTimeout(context.Background(), WriteTimeout)
	defer cancel()

	kafkaMessages, err := mapEvents(events)
	if err != nil {
		b.log.Errorf("failed to map events: %v", err)
		return
	}
	err = b.writer.Write(ctx, kafkaMessages)
	if err != nil {
		// TODO: Add publishing to DLQ
		metrics.PublishErrors.Inc()
		b.log.Errorf("failed to publish events: %v", err)
	}
}

func mapEvents(events []models.Event) ([]kafkago.Message, error) {
	result := make([]kafkago.Message, 0, len(events))
	for _, event := range events {
		raw, err := json.Marshal(kafkamodels.TelemetryEvent{
			DeviceID:  event.DeviceID,
			Type:      string(event.Type),
			Value:     event.Value,
			Timestamp: event.Timestamp,
		})
		if err != nil {
			return nil, err
		}
		result = append(result, kafkago.Message{
			Key:   []byte(event.DeviceID),
			Value: raw,
		})
	}
	return result, nil
}
