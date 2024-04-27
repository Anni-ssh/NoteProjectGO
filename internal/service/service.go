package service

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/storage"
)

type Authorization interface {
	CreateUser(user entities.User) (int, error)
	CheckUser(username, password string) (*entities.User, error)
	GenAuthToken(user entities.User) (string, error)
}

type Service struct {
	Authorization
}

func NewService(r *storage.Storage) *Service {
	return &Service{
		Authorization: NewAuthService(r.Authorization),
	}
}
