package handler

import (
	"NoteProject/internal/entities"
	"NoteProject/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type userInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// @Summary signUp Регистрация в приложенеи
// @Description Обрабатывает запрос на регистрацию в приложение.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body userInput true "Данные для создания пользователя"
// @Header 200 {string} Authorization "Bearer"
// @Failure 400 {object} Response "Ошибка валидации"
// @Failure 401 {object} Response "Ошибка создания заметки"
// @Failure 422 {object} Response "Неправильный формат данных"
// @Failure 500 {object} Response "Внутренняя ошибка сервера"
// @Router /auth/sign-up [post]
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

		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		msg := ErrValidator(validateErr)

		NewErrResponse(w, http.StatusUnprocessableEntity, msg.ErrMsg)
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

// @Summary signIn Вход в приложение
// @Description Обрабатывает запрос предназначенный для входа в приложение.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body userInput true "Данные для создания заметки"
// @Header 200 {string} Authorization "Bearer"
// @Failure 400 {object} Response "Ошибка валидации"
// @Failure 401 {object} Response "Ошибка создания заметки"
// @Failure 422 {object} Response "Неправильный формат данных"
// @Failure 500 {object} Response "Внутренняя ошибка сервера"
// @Router /auth/sign-in [post]
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = h.services.Session.CreateSession(ctx, strconv.Itoa(regUser.Id), token, time.Hour*24)
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
