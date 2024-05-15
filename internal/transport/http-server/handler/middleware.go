package handler

import (
	"NoteProject/pkg/logger"
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const authHeader = "Authorization"

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.authMiddleware"
		log := h.Logs.With(slog.String("operation", op))

		header := r.Header.Get(authHeader)

		if header == "" {
			NewErrResponse(w, http.StatusUnauthorized, "Authorization Token is empty")
			log.Error("userToken is empty")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			NewErrResponse(w, http.StatusUnauthorized, "Authorization Token is invalid")
			log.Error("userToken is invalid")
			return
		}

		token := headerParts[1]
		userID, err := h.services.Authorization.ParseToken(token)
		if err != nil {
			NewErrResponse(w, http.StatusUnauthorized, "Authorization Token is invalid")
			log.Error("userToken is invalid", logger.Err(err))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = h.services.Session.CheckSession(ctx, strconv.Itoa(userID), token)
		if err != nil {
			NewErrResponse(w, http.StatusUnauthorized, "Authorization Token is invalid")
			log.Error("userToken is invalid", logger.Err(err))
			return
		}

		next.ServeHTTP(w, r)
	})
}
