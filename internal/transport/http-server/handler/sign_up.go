package handler

import (
	"NoteProject/internal/entities"
	"NoteProject/pkg/logger"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	const op = "handler.signUp"
	log := h.Logs.With(slog.String("operation", op))

	var user entities.User

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

	if err := validator.New().Struct(user); err != nil {
		log.Error("Failed to validate data", logger.Err(err))

		validateErr := err.(validator.ValidationErrors)
		msg := ErrValidator(validateErr)

		NewErrResponse(w, http.StatusUnprocessableEntity, msg.ErrorMsg)
		return
	}

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		NewErrResponse(w, http.StatusUnauthorized, "Create user failed")
		log.Error("Create user failed", logger.Err(err))
		return
	}

	response := map[string]interface{}{
		"id": id,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Server Error: error encode JSON response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server Error")
		return
	}
	log.Info("Successfully created user", "user", user.Username)
}
