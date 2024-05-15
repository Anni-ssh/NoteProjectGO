package handler

import (
	"log/slog"
	"net/http"
)

func (h *Handler) noteWorkspace(w http.ResponseWriter, r *http.Request) {
	const op = "handler.noteWorkspace"
	log := h.Logs.With(slog.String("operation", op))

	_, _ = w.Write([]byte("It is workspace"))
	log.Info("It is workspace")

}
