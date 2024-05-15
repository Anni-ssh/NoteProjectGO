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

type inputNote struct {
	UserId int    `json:"user_id" validate:"required"`
	Title  string `json:"title" validate:"required"`
	Text   string `json:"text" validate:"required"`
}

type reqNote struct {
	ID int `json:"note_id" validate:"required"`
}

type reqUser struct {
	ID int `json:"user_id" validate:"required"`
}

const (
	errEmptyBody   = "Request body is empty"
	errInvalidBody = "Request body is invalid"
	errInvalidData = "Request data is invalid"
)

// @Summary Create Создать заметку
// @Description Обрабатывает запрос на создание новой заметки.
// @Tags Notes
// @Accept json
// @Produce json
// @Param body body inputNote true "Данные для создания заметки"
// @Success 200 {object} map[string]interface{} "Успешный ответ"
// @Failure 400 {object} Response "Ошибка валидации"
// @Failure 401 {object} Response "Ошибка создания заметки"
// @Failure 422 {object} Response "Неправильный формат данных"
// @Failure 500 {object} Response "Внутренняя ошибка сервера"
// @Router /note/create [post]
// @Security ApiKeyAuth
func (h *Handler) noteCreate(w http.ResponseWriter, r *http.Request) {
	const op = "handler.noteCreate"
	log := h.Logs.With(slog.String("operation", op))

	var note inputNote

	err := render.DecodeJSON(r.Body, &note)
	if errors.Is(err, io.EOF) {
		log.Error(errEmptyBody, logger.Err(err))
		NewErrResponse(w, http.StatusUnprocessableEntity, errEmptyBody)
		return
	}

	if err != nil {
		log.Error(errEmptyBody, logger.Err(err))
		NewErrResponse(w, http.StatusUnprocessableEntity, errInvalidBody)
		return
	}

	if err := validator.New().Struct(note); err != nil {

		validateErr := err.(validator.ValidationErrors)
		msg := ErrValidator(validateErr)

		log.Error(msg.ErrMsg, logger.Err(err))

		NewErrResponse(w, http.StatusUnprocessableEntity, msg.ErrMsg)
		return
	}

	id, err := h.services.Note.CreateNote(note.UserId, note.Title, note.Text)
	if err != nil {
		NewErrResponse(w, http.StatusInternalServerError, "Create note failed")
		log.Error("Create note failed", logger.Err(err))
		return
	}

	response := map[string]interface{}{
		"id": id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Server Error: error encode JSON response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server Error")
		return
	}
	log.Info("Successfully created note", "note", note.Title)
}

// @Summary List Показать заметки
// @Description Обрабатывает запрос на отображение заметок.
// @Tags Notes
// @Accept json
// @Produce json
// @Param body body reqUser true "Данные для отображения заметок"
// @Success 200 {object} map[string]interface{} "Успешный ответ"
// @Failure 400 {object} Response "Ошибка валидации"
// @Failure 401 {object} Response "Ошибка создания заметки"
// @Failure 422 {object} Response "Неправильный формат данных"
// @Failure 500 {object} Response "Внутренняя ошибка сервера"
// @Router /note/list [post]
// @Security ApiKeyAuth
func (h *Handler) notesList(w http.ResponseWriter, r *http.Request) {
	const op = "handler.notesList"
	log := h.Logs.With(slog.String("operation", op))

	var user reqUser

	err := render.DecodeJSON(r.Body, &user)
	if errors.Is(err, io.EOF) {
		log.Error(errEmptyBody, logger.Err(err))
		NewErrResponse(w, http.StatusUnprocessableEntity, errEmptyBody)
		return
	}
	if err != nil {
		log.Error(errInvalidBody, logger.Err(err))
		NewErrResponse(w, http.StatusUnprocessableEntity, errInvalidBody)
		return
	}

	notesList, err := h.services.Note.NotesList(user.ID)
	if err != nil {
		NewErrResponse(w, http.StatusBadRequest, "Getting notes list failed")
		log.Error("Getting notes list failed", logger.Err(err))
		return
	}

	response := map[string]interface{}{
		"notes": notesList,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Server Error: error encode JSON response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server Error")
		return
	}
	log.Info("Notes list successfully sent to user", "userID", user.ID)
}

// @Summary Update Обновить заметку
// @Description Обрабатывает запрос на обновлние заметоки.
// @Tags Notes
// @Accept json
// @Produce json
// @Param body body entities.Note true "Данные для обновления заметки"
// @Success 200 {object} map[string]interface{} "Успешный ответ"
// @Failure 400 {object} Response "Ошибка валидации"
// @Failure 401 {object} Response "Ошибка создания заметки"
// @Failure 422 {object} Response "Неправильный формат данных"
// @Failure 500 {object} Response "Внутренняя ошибка сервера"
// @Router /note/update [put]
// @Security ApiKeyAuth
func (h *Handler) noteUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "handler.noteUpdate"
	log := h.Logs.With(slog.String("operation", op))

	var note entities.Note

	err := render.DecodeJSON(r.Body, &note)
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

	err = h.services.Note.UpdateNote(note)
	if err != nil {
		NewErrResponse(w, http.StatusBadRequest, "Update note failed")
		log.Error("Update note failed", logger.Err(err))
		return
	}

	response := map[string]interface{}{
		"Response": "Successfully updated note",
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Server Error: error encode JSON response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server Error")
		return
	}
	log.Info("Notes successfully update", "noteID", note.Id)
}

// @Summary Update Удалить заметку
// @Description Обрабатывает запрос на удаление заметок.
// @Tags Notes
// @Accept json
// @Produce json
// @Param body body reqNote true "Данные для удаления заметки"
// @Success 200 {object} map[string]interface{} "Успешный ответ"
// @Failure 400 {object} Response "Ошибка валидации"
// @Failure 401 {object} Response "Ошибка создания заметки"
// @Failure 422 {object} Response "Неправильный формат данных"
// @Failure 500 {object} Response "Внутренняя ошибка сервера"
// @Router /note/delete [delete]
// @Security ApiKeyAuth
func (h *Handler) noteDelete(w http.ResponseWriter, r *http.Request) {
	const op = "handler.noteDelete"
	log := h.Logs.With(slog.String("operation", op))

	var note reqNote

	err := render.DecodeJSON(r.Body, &note)
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

	err = h.services.Note.DeleteNote(note.ID)
	if err != nil {
		NewErrResponse(w, http.StatusBadRequest, "Delete note failed")
		log.Error("Delete note failed", logger.Err(err))
		return
	}

	response := map[string]interface{}{
		"Response": "Note successfully delete",
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Server Error: error encode JSON response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Server Error")
		return
	}
	log.Info("Note successfully delete", "noteID", note.ID)
}
