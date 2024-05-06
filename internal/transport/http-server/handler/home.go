package handler

import (
	"NoteProject/internal/transport/http-server/paths"
	"NoteProject/pkg/logger"
	"html/template"
	"log/slog"
	"net/http"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	const op = "handler.home"
	log := h.Logs.With(slog.String("operation", op))

	tmpl, err := template.ParseFiles(paths.TemplatesPath(paths.Home))
	if err != nil {
		log.Error("failed to parse tmpl files", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server error")
		return
	}

	if err = tmpl.Execute(w, nil); err != nil {
		log.Error("failed to execute tmpl", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server error")
		return
	}
}
