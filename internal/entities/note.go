package entities

import "time"

// Note represents the structure of a note
// @Description Note structure
// @Accept  json
// @Produce  json
// @Param id query int true "Note ID"
// @Param user_id query int true "User ID"
// @Param title query string true "Title of the note"
// @Param text query string true "Text of the note"
// @Param date query string true "Date of the note in RFC3339 format"
// @Param done query bool true "Completion status of the note"
// @Success 200 {object} Note "Successfully retrieved note"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Router /notes [post]
type Note struct {
	Id     int       `json:"id"`
	UserId int       `json:"user_id" validate:"required"`
	Title  string    `json:"title" validate:"required"`
	Text   string    `json:"text" validate:"required"`
	Date   time.Time `json:"date"`
	Done   bool      `json:"done"`
}
