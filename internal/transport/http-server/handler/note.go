package handler

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/errs"
	"NoteProject/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type inputNote struct {
	UserId int    `json:"user_id"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

// @Summary Create a Note
// @Description Handles the request to create a new note.
// @Tags Notes
// @Accept json
// @Produce json
// @Param body body inputNote true "Data for creating a note"
// @Success 201 {object} map[string]interface{} "Return id"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal server error"
// @Router /note/create [post]
// @Security ApiKeyAuth
func (h *Handler) noteCreate(w http.ResponseWriter, r *http.Request) {
	const op = "handler.noteCreate"
	log := h.Logs.With(slog.String("operation", op))

	var note inputNote

	// Decode JSON request body
	err := render.DecodeJSON(r.Body, &note)
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

	// Create note
	id, err := h.services.Note.CreateNote(note.UserId, note.Title, note.Text)
	if err != nil {

		if errors.Is(err, errs.ErrUserNotExists) {
			log.Error("invalid user id", logger.Err(err))
			NewErrResponse(w, http.StatusBadRequest, "Invalid user id")
			return
		}

		log.Error("Create note failed", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Response
	response := map[string]interface{}{
		"id": id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Server Error: error encode JSON response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}

// @Summary List Show users notes
// @Description Handles the request to display notes.
// @Tags Notes
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {array} entities.Note
// @Failure 400 {object} Response "Bad request"
// @Failure 404 {object} Response "Note not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /note/list/{userID} [get]
// @Security ApiKeyAuth
func (h *Handler) notesList(w http.ResponseWriter, r *http.Request) {
	const op = "handler.notesList"
	log := h.Logs.With(slog.String("operation", op))

	userID := chi.URLParam(r, "userID")
	id, err := strconv.Atoi(userID)
	if err != nil {
		NewErrResponse(w, http.StatusBadRequest, "Invalid userID")
		log.Error("Invalid userID", logger.Err(err))
		return
	}

	notesList, err := h.services.Note.NotesList(id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotExists) {
			log.Error("Invalid user id", logger.Err(err))
			NewErrResponse(w, http.StatusNotFound, "User does not exist")
			return
		}

		NewErrResponse(w, http.StatusInternalServerError, "Internal server error")
		log.Error("Getting notes list failed", logger.Err(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(notesList) == 0 {
		response := map[string]interface{}{
			"Response": fmt.Sprintf("User %d does not have notes", id),
		}
		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Error("Server Error: error encoding JSON response", logger.Err(err))
			NewErrResponse(w, http.StatusInternalServerError, "Intermal server Error")
			return
		}
		return
	}

	if err = json.NewEncoder(w).Encode(notesList); err != nil {
		log.Error("Server Error: error encoding JSON response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Intermal server Error")
		return
	}
}

// @Summary Update the note
// @Description Handles the request to update a note.
// @Tags Notes
// @Accept json
// @Produce json
// @Param body body entities.Note true "Data for updating the note"
// @Success 200 {object} map[string]interface{} "Successfully updated note"
// @Failure 400 {object} Response "Bad request"
// @Failure 404 {object} Response "Note not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /note/update [put]
// @Security ApiKeyAuth
func (h *Handler) noteUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "handler.noteUpdate"
	log := h.Logs.With(slog.String("operation", op))

	var note entities.Note

	// Decode JSON request body
	err := render.DecodeJSON(r.Body, &note)
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

	err = h.services.Note.UpdateNote(note)
	if err != nil {
		if errors.Is(err, errs.ErrNoteNotExists) {
			log.Error("Note does not exist", logger.Err(err))
			NewErrResponse(w, http.StatusNotFound, "Note does not exist")
			return
		}

		log.Error("Failed to update note", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := map[string]interface{}{
		"Response": "Successfully updated note",
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Server Error: error encoding JSON response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Intermal server Error")
		return
	}
}

// @Summary Delete the note
// @Description Handles the request to delete a note.
// @Tags Notes
// @Accept json
// @Produce json
// @Param noteID path int true "Note ID"
// @Success 200 {object} map[string]interface{} "Successful response"
// @Failure 400 {object} Response "Bad request"
// @Failure 404 {object} Response "Note not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /note/delete/{noteID} [delete]
// @Security ApiKeyAuth
func (h *Handler) noteDelete(w http.ResponseWriter, r *http.Request) {
	const op = "handler.noteDelete"
	log := h.Logs.With(slog.String("operation", op))

	noteID := chi.URLParam(r, "noteID")
	id, err := strconv.Atoi(noteID)
	if err != nil {
		NewErrResponse(w, http.StatusBadRequest, "Invalid noteID")
		log.Error("Invalid noteID", logger.Err(err))
		return
	}

	err = h.services.Note.DeleteNote(id)

	if err != nil {
		if errors.Is(err, errs.ErrNoteNotExists) {
			log.Error("Note does not exist", logger.Err(err))
			NewErrResponse(w, http.StatusNotFound, "Note does not exist")
			return
		}

		log.Error("Delete note failed", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := map[string]interface{}{
		"Response": "Note successfully deleted",
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Server Error: error encoding JSON response", logger.Err(err))
		NewErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}
