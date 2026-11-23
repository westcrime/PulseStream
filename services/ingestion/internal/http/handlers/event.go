package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/westcrime77/PulseStream/services/ingestion/internal/http/models"
	"github.com/westcrime77/PulseStream/services/ingestion/internal/metrics"
)

func (h *Handler) postEvents(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(metrics.PublishDuration)
	defer timer.ObserveDuration()

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

	for _, e := range events {
		metrics.EventsReceived.WithLabelValues(string(e.Type)).Inc()
	}

	if rest := h.batcher.Add(events...); len(rest) > 0 {
		for range rest {
			metrics.EventsPublishedWithoutBatcher.Inc()
		}
		h.batcher.PublishEvents(rest)
	}
	w.WriteHeader(http.StatusAccepted)
}

func validateEvents(events []models.Event) error {
	for _, event := range events {
		if err := event.Validate(); err != nil {
			return err
		}
	}
	return nil
}
