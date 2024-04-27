package handler

import (
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.userIdentity"
		log := h.Logs.With(slog.String("operation", op))

		header := r.Header.Get("Authorization")

		if header == "" {
			NewErrResponse(w, http.StatusUnauthorized, "userToken is empty")
			log.Error("userToken is empty")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			NewErrResponse(w, http.StatusUnauthorized, "userToken is invalid")
			log.Error("userToken is invalid")
			return
		}

		//token := headerParts[1]

		next.ServeHTTP(w, r)
	})
}
