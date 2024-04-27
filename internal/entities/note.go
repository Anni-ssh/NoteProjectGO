package entities

import "time"

type Note struct {
	Id     int       `json:"id" validate:"required"`
	UserId int       `json:"userId" validate:"required"`
	Title  string    `json:"title" validate:"required"`
	Text   string    `json:"text" validate:"required"`
	Date   time.Time `json:"date" validate:"required"`
	Done   bool      `json:"done" validate:"required"`
}
