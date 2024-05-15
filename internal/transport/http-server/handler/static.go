package handler

import (
	"net/http"
)

func (h *Handler) static(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("static method"))
}
