package handler

import "net/http"

func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete method"))
}
