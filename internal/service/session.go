package service

import "NoteProject/internal/storage"

type SessionService struct {
	storage storage.Session
}

func NewSessionService(storage storage.Session) *SessionService {
	return &SessionService{storage: storage}
}

func (s *SessionService) CreateSession(userID, token string) error {
	return s.storage.CreateSession(userID, token)
}

func (s *SessionService) CheckSession(userID, token string) error {
	return s.storage.CheckSession(userID, token)
}
