package handler

import "net/http"

func (h *Handler) saveNote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("save method"))
}
