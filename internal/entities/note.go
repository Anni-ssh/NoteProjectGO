package entities

import "time"

type Note struct {
	Id     int       `json:"id"`
	UserId int       `json:"user_id" validate:"required"`
	Title  string    `json:"title" validate:"required"`
	Text   string    `json:"text" validate:"required"`
	Date   time.Time `json:"date"`
	Done   bool      `json:"done"`
}
