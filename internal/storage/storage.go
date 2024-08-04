package storage

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/storage/postgres"
	"NoteProject/internal/storage/redisDB"
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=storage.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user entities.User) (int, error)
	CheckUser(username string, password string) (entities.User, error)
}
type Session interface {
	CreateSession(ctx context.Context, userID, token string, expiration time.Duration) error
	CheckSession(ctx context.Context, userID, token string) error
}

type NoteManage interface {
	CreateNote(userID int, title, text string) (int, error)
	NotesList(userID int) ([]entities.Note, error)
	UpdateNote(note entities.Note) error
	DeleteNote(noteID int) error
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
