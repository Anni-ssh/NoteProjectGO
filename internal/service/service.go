package service

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/storage"
)

type Authorization interface {
	CreateUser(user entities.User) (int, error)
	CheckUser(username, password string) (*entities.User, error)
	GenToken(user entities.User) (string, error)
	ParseToken(accessToken string) (int, error)
}

type Session interface {
	CreateSession(userID, token string) error
	CheckSession(userID, token string) error
}

type Service struct {
	Authorization
	Session
}

func NewService(r *storage.Storage) *Service {
	return &Service{
		Authorization: NewAuthService(r.Authorization),
		Session:       NewSessionService(r.Session),
	}
}
