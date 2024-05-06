package handler

import (
	"NoteProject/pkg/logger"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
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

	regUser, err := h.services.Authorization.CheckUser(user.Username, user.Password)
	if err != nil {
		log.Error("invalid user data", logger.Err(err))
		NewErrResponse(w, http.StatusUnauthorized, "Invalid user data")
		return
	}

	token, err := h.services.Authorization.GenToken(*regUser)
	if err != nil {
		log.Error("failed to create a token", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server error")
		return
	}

	err = h.services.Session.CreateSession(strconv.Itoa(regUser.Id), token)
	if err != nil {
		log.Error("failed to create session", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server error")
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Login successful"))
	if err != nil {
		log.Error("failed to write response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server Error")
		return
	}
}

func (h *Handler) Session(w http.ResponseWriter, r *http.Request) {

	var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	const op = "handler.Session"
	log := h.Logs.With(slog.String("operation", op))

	session, _ := store.Get(r, "session-note-project")

	session.Values[42] = "dsfsafwefw"
	err := session.Save(r, w)
	if err != nil {
		log.Error("failed to write response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server Error")
		return
	}
}
