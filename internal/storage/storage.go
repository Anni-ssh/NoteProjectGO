package storage

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/storage/postgres"
	"NoteProject/internal/storage/redisDB"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const (
	usersTable = "users"
	notesTable = "notes"
)

type Authorization interface {
	CreateUser(user entities.User) (int, error)
	CheckUser(username string, password string) (*entities.User, error)
}
type Session interface {
	CreateSession(userID, token string) error
	CheckSession(userID, token string) error
}

type NoteManage interface {
	CreateNote() error
	DeleteNote() error
	UpdateNote() error
}

type Storage struct {
	Authorization
	Session
	NoteManage
}

func NewStorage(db *sql.DB, r *redis.Client) *Storage {
	return &Storage{
		Authorization: postgres.NewAuthPostgres(db),
		Session:       redisDB.NewRedisStorage(r),
		NoteManage:    postgres.NewNoteManage(db),
	}
}
