package service

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/storage"
	"context"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user entities.User) (int, error)
	CheckUser(username, password string) (entities.User, error)
	GenToken(user entities.User) (string, error)
	ParseToken(accessToken string) (int, error)
}

type Session interface {
	CreateSession(ctx context.Context, userID, token string, expiration time.Duration) error
	CheckSession(ctx context.Context, userID, token string) error
}

type Note interface {
	CreateNote(userID int, title, text string) (int, error)
	NotesList(userID int) ([]entities.Note, error)
	UpdateNote(note entities.Note) error
	DeleteNote(noteID int) error
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
