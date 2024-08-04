package handler

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/errs"
	"NoteProject/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type userInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// @Summary signUp User registration
// @Description Handles the request to register a user in the application.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body userInput true "User data for creation"
// @Success 201 {object} map[string]interface{} "User ID"
// @Failure 400 {object} Response "Bad request"
// @Failure 409 {object} Response "Conflict: User already exists"
// @Failure 500 {object} Response "Internal server error"
// @Router /auth/sign-up [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	const op = "handler.signUp"
	log := h.Logs.With(slog.String("operation", op))

	var user entities.User
	// Parse JSON
	err := render.DecodeJSON(r.Body, &user)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("Request body is empty", logger.Err(err))
			NewErrResponse(w, http.StatusBadRequest, "Request body is empty")
			return
		}

		log.Error("Failed to decode request body", logger.Err(err))
		NewErrResponse(w, http.StatusBadRequest, "Invalid JSON request")
		return
	}

	// Validate data
	if err := validator.New().Struct(user); err != nil {
		log.Error("Failed to validate data", logger.Err(err))

		var validateErr validator.ValidationErrors
		if errors.As(err, &validateErr) {
			msg := ErrValidator(validateErr)
			NewErrResponse(w, http.StatusBadRequest, msg.ErrMsg)
			return
		}

		log.Error("Validation error", logger.Err(err))
		NewErrResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}
	// Context
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Create user
	id, err := h.services.Authorization.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, errs.ErrUserExists) {
			log.Error("Create user failed: user already exists", logger.Err(err))
			NewErrResponse(w, http.StatusConflict, "Username already exists")
			return
		}

		log.Error("Create user failed", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Уведомление о регистрации пользователя
	log.Info("Sucsessful registration", "User", user.Username)

	// Prepare and send response
	response := map[string]interface{}{
		"id": id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Server error: failed to encode JSON response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// @Summary signIn Log in to the app
// @Description Handles the request intended for logging into the application.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body userInput true "User data for login"
// @Success 200 {string} string "Login successful"
// @Failure 400 {object} Response "Bad request"
// @Failure 401 {object} Response "Unauthorized access"
// @Failure 500 {object} Response "Internal server error"
// @Router /auth/sign-in [post]
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	const op = "handler.signIn"
	log := h.Logs.With(slog.String("operation", op))

	var user userInput

	// Decode JSON request body
	err := render.DecodeJSON(r.Body, &user)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", logger.Err(err))
			NewErrResponse(w, http.StatusBadRequest, "Request body is empty")
			return
		}
		log.Error("failed to decode request body", logger.Err(err))
		NewErrResponse(w, http.StatusBadRequest, "Invalid JSON request")
		return
	}

	// Validate user input data
	if err := validator.New().Struct(user); err != nil {
		log.Error("failed to validate user data", logger.Err(err))

		var validateErr validator.ValidationErrors
		if errors.As(err, &validateErr) {
			msg := ErrValidator(validateErr)
			NewErrResponse(w, http.StatusBadRequest, msg.ErrMsg)
			return
		}

		NewErrResponse(w, http.StatusBadRequest, "Invalid data")
		return
	}

	// Context
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check user
	regUser, err := h.services.Authorization.CheckUser(ctx, user.Username, user.Password)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotExists) {
			log.Error("invalid user credentials", logger.Err(err))
			NewErrResponse(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		log.Error("failed to check user credentials", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Generate token
	token, err := h.services.Authorization.GenToken(regUser)
	if err != nil {
		log.Error("failed to generate token", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Create session
	err = h.services.Session.CreateSession(r.Context(), strconv.Itoa(regUser.Id), token, 24*time.Hour)
	if err != nil {
		log.Error("failed to create session", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Уведомление о входе пользователя
	log.Info("Sucsessful login", "User", user.Username)

	// Set authorization header and respond
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Login successful"))
	if err != nil {
		log.Error("failed to write response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}
