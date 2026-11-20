package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/westcrime77/PulseStream/pkg/kafka"
)

type Logger interface {
	Errorf(format string, v ...any)
}

type Handler struct {
	writer *kafka.Writer
	log    Logger
}

func New(writer *kafka.Writer) *Handler {
	return &Handler{
		writer: writer,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Get("/health", h.health)
	r.Post("/events", h.postEvents)

	return r
}
