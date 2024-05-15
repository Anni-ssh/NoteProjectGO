package service

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/storage"
)

type NoteService struct {
	storage storage.NoteManage
}

func NewNoteManage(s storage.NoteManage) *NoteService {
	return &NoteService{storage: s}
}

func (n *NoteService) CreateNote(userID int, title, text string) (int, error) {
	return n.storage.CreateNote(userID, title, text)
}

func (n *NoteService) NotesList(userID int) ([]entities.Note, error) {
	return n.storage.NotesList(userID)
}

func (n *NoteService) DeleteNote(noteID int) error {
	return n.storage.DeleteNote(noteID)
}

func (n *NoteService) UpdateNote(note entities.Note) error {
	return n.storage.UpdateNote(note)
}
