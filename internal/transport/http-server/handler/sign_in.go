package handler

import (
	"NoteProject/pkg/logger"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

type userInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	const op = "handler.signIn"
	log := h.Logs.With(slog.String("operation", op))

	var user userInput

	err := render.DecodeJSON(r.Body, &user)
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty", logger.Err(err))
		NewErrResponse(w, http.StatusUnprocessableEntity, "Request body is empty")
		return
	}

	if err != nil {
		log.Error("failed to decode request body", logger.Err(err))
		NewErrResponse(w, http.StatusUnprocessableEntity, "Failed to decode request")
		return
	}
	if err = validator.New().Struct(user); err != nil {
		log.Error("Failed to validate data", logger.Err(err))
		NewErrResponse(w, http.StatusUnprocessableEntity, "Failed to validate data")
		return
	}

	result, err := h.services.Authorization.CheckUser(user.Username, user.Password)
	if err != nil {
		NewErrResponse(w, http.StatusUnauthorized, "invalid user data")
		log.Error("invalid user data", logger.Err(err))
		return
	}

	token, err := h.services.Authorization.GenAuthToken(*result)
	if err != nil {
		NewErrResponse(w, http.StatusInternalServerError, "server error")
		log.Error("failed to create a token", logger.Err(err))
		return
	}

	response := map[string]interface{}{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Server Error: error encode JSON response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server Error")
		return
	}
}
