package handlers

import (
	"net/http"
)

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}
