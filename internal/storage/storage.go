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
	CreateUser(ctx context.Context, user entities.User) (int, error)
	CheckUser(ctx context.Context, username string, password string) (entities.User, error)
}
type Session interface {
	CreateSession(ctx context.Context, userID, token string, expiration time.Duration) error
	CheckSession(ctx context.Context, userID, token string) error
}

type NoteManage interface {
	CreateNote(ctx context.Context, userID int, title, text string) (int, error)
	NotesList(ctx context.Context, userID int) ([]entities.Note, error)
	UpdateNote(ctx context.Context, note entities.Note) error
	DeleteNote(ctx context.Context, noteID int) error
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
