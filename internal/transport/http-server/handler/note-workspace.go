package handler

import (
	"net/http"
)

func (h *Handler) noteWorkspace(w http.ResponseWriter, r *http.Request) {
	const op = "handler.noteWorkspace"

	_, _ = w.Write([]byte("It is workspace"))

}
