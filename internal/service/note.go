package service

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/storage"
	"context"
)

type NoteService struct {
	storage storage.NoteManage
}

func NewNoteManage(s storage.NoteManage) *NoteService {
	return &NoteService{storage: s}
}

func (n *NoteService) CreateNote(ctx context.Context, userID int, title, text string) (int, error) {
	return n.storage.CreateNote(ctx, userID, title, text)
}

func (n *NoteService) NotesList(ctx context.Context, userID int) ([]entities.Note, error) {
	return n.storage.NotesList(ctx, userID)
}

func (n *NoteService) DeleteNote(ctx context.Context, noteID int) error {
	return n.storage.DeleteNote(ctx, noteID)
}

func (n *NoteService) UpdateNote(ctx context.Context, note entities.Note) error {
	return n.storage.UpdateNote(ctx, note)
}
