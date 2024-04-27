package handler

import "net/http"

func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create method"))
}
