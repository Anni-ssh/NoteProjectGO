package service

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/storage"
	"context"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(ctx context.Context, user entities.User) (int, error)
	CheckUser(ctx context.Context, username, password string) (entities.User, error)
	GenToken(user entities.User) (string, error)
	ParseToken(accessToken string) (int, error)
}

type Session interface {
	CreateSession(ctx context.Context, userID, token string, expiration time.Duration) error
	CheckSession(ctx context.Context, userID, token string) error
}

type Note interface {
	CreateNote(ctx context.Context, userID int, title, text string) (int, error)
	NotesList(ctx context.Context, userID int) ([]entities.Note, error)
	UpdateNote(ctx context.Context, note entities.Note) error
	DeleteNote(ctx context.Context, noteID int) error
}

type Service struct {
	Authorization
	Session
	Note
}

func NewService(s *storage.Storage) *Service {
	return &Service{
		Authorization: NewAuthService(s.Authorization),
		Session:       NewSessionService(s.Session),
		Note:          NewNoteManage(s.NoteManage),
	}
}
