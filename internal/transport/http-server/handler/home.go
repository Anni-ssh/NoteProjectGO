package handler

import "net/http"

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("IT IS HOME PAGE"))
}
