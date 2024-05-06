package handler

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/transport/http-server/paths"
	"NoteProject/pkg/logger"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) noteWorkspace(w http.ResponseWriter, r *http.Request) {
	const op = "handler.noteWorkspace"
	log := h.Logs.With(slog.String("operation", op))

	tmpl, err := template.ParseFiles(paths.TemplatesPath(paths.NoteWorkSpace), paths.TemplatesPath(paths.Note))
	if err != nil {
		log.Error("failed to parse tmpl files", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server error")
		return
	}

	var notes []entities.Note

	note := entities.Note{
		Id:     50,
		UserId: 50,
		Title:  "Заметка",
		Text:   "Что-то важное",
		Date:   time.Now(),
		Done:   false,
	}

	notes = append(notes, note)

	if err = tmpl.Execute(w, notes); err != nil {
		log.Error("failed to execute tmpl", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server error")
		return
	}

}
