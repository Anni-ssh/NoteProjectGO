package handler

import (
	"log/slog"
	"net/http"
)

func (h *Handler) static(w http.ResponseWriter, r *http.Request) {
	slog.Info("handler", "Static")
}
