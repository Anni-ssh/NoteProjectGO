package storage

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/storage/postgres"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	usersTable = "users"
	notesTable = "notes"
)

type Authorization interface {
	CreateUser(user entities.User) (int, error)
	CheckUser(username string, password string) (*entities.User, error)
}

type Storage struct {
	Authorization
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Authorization: postgres.NewAuthPostgres(db),
	}
}
