package postgres

import (
	"NoteProject/internal/entities"
	"database/sql"
	"fmt"
)

type NoteManagePostgres struct {
	db *sql.DB
}

func NewNoteManage(db *sql.DB) *NoteManagePostgres {
	return &NoteManagePostgres{db: db}
}

func (n *NoteManagePostgres) CreateNote(userID int, title, text string) (int, error) {
	const operation = "postgres.CreateNote"
	q, err := n.db.Prepare("INSERT INTO notes (user_id, title, text) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s Prepare: %w", operation, err)
	}

	var id int

	err = q.QueryRow(userID, title, text).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s Scan: %w", operation, err)
	}

	return id, nil
}

func (n *NoteManagePostgres) GetNotesList(userID int) ([]entities.Note, error) {
	const operation = "postgres.GetNotesList"
	q, err := n.db.Prepare("SELECT * FROM notes WHERE user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s Prepare: %w", operation, err)
	}

	var notes []entities.Note

	rows, err := q.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("%s Query: %w", operation, err)
	}

	for rows.Next() {
		var note entities.Note
		err := rows.Scan(&note.Id, &note.UserId, &note.Title, &note.Text, &note.Date, &note.Done)
		if err != nil {
			return nil, fmt.Errorf("%s Scan: %w", operation, err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s Rows: %w", operation, err)
	}

	return notes, nil

}

func (n *NoteManagePostgres) DeleteNote(noteID int) error {
	const operation = "postgres.DeleteNote"
	q, err := n.db.Prepare("DELETE FROM notes WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s Prepare: %w", operation, err)
	}

	_, err = q.Exec(noteID)
	if err != nil {
		return fmt.Errorf("%s Exec: %w", operation, err)
	}

	return nil
}

func (n *NoteManagePostgres) UpdateNote(userID, title, text string) (int, error) {
	return 0, nil
}
