package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/westcrime77/PulseStream/pkg/kafka"
	kafkamodels "github.com/westcrime77/PulseStream/pkg/kafka/models"
	"github.com/westcrime77/PulseStream/services/ingestion/internal/http/models"
)

func (h *Handler) postEvents(w http.ResponseWriter, r *http.Request) {
	var events []models.Event
	err := json.NewDecoder(r.Body).Decode(&events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Print("Error while closing the body",
				err.Error())
		}
	}(r.Body)

	if err = validateEvents(events); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), 10*time.Second)
	go func() {
		defer cancel()
		if err := publishEvents(ctx, h.writer, events); err != nil {
			h.log.Errorf("publish events: %v", err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}

func publishEvents(ctx context.Context, writer *kafka.Writer, events []models.Event) error {
	kafkaMessages, err := mapEvents(events)
	if err != nil {
		return err
	}
	return writer.Write(ctx, kafkaMessages)
}

func mapEvents(events []models.Event) ([]kafkago.Message, error) {
	result := make([]kafkago.Message, 0, len(events))
	for _, event := range events {
		raw, err := json.Marshal(kafkamodels.TelemetryEvent{
			DeviceID: event.DeviceID,
			Type:     string(event.Type),
			Value:    event.Value,
		})
		if err != nil {
			return nil, err
		}
		result = append(result, kafkago.Message{
			Key:   nil,
			Value: raw,
		})
	}
	return result, nil
}

func validateEvents(events []models.Event) error {
	for _, event := range events {
		if err := event.Validate(); err != nil {
			return err
		}
	}
	return nil
}
