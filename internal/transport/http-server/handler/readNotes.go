package handler

import "net/http"

func (h *Handler) readNote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("read method"))
}
