package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/westcrime77/PulseStream/services/ingestion/internal/batcher"
)

type Logger interface {
	Errorf(format string, v ...any)
}

type Handler struct {
	log     Logger
	batcher *batcher.Batcher
}

func New(batcher *batcher.Batcher, log Logger) *Handler {
	return &Handler{
		log:     log,
		batcher: batcher,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Get("/health", h.health)
	r.Handle("/metrics", promhttp.Handler())
	r.Post("/events", h.postEvents)

	return r
}
