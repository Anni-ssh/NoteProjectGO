package service

import (
	"NoteProject/internal/storage"
	"context"
	"time"
)

type SessionService struct {
	storage storage.Session
}

func NewSessionService(storage storage.Session) *SessionService {
	return &SessionService{storage: storage}
}

func (s *SessionService) CreateSession(ctx context.Context, userID, token string, expiration time.Duration) error {
	return s.storage.CreateSession(ctx, userID, token, expiration)
}

func (s *SessionService) CheckSession(ctx context.Context, userID, token string) error {
	return s.storage.CheckSession(ctx, userID, token)
}
