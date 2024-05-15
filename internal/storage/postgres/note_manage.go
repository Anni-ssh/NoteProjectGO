package postgres

import (
	"NoteProject/internal/entities"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
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

func (n *NoteManagePostgres) NotesList(userID int) ([]entities.Note, error) {
	const operation = "postgres.NotesList"
	q, err := n.db.Prepare("SELECT * FROM notes WHERE user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s Prepare: %w", operation, err)
	}

	var notes []entities.Note

	rows, err := q.Query(userID)
	if err != nil {
		var pgErr *pq.Error
		ok := errors.As(err, &pgErr)
		if !ok {
			return nil, fmt.Errorf("%s Query: %w", operation, err)
		}

		if pgErr.Code.Name() == uniqueViolationCode {
			return nil, errNoteNotFound
		}

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

	result, err := q.Exec(noteID)
	if err != nil {
		return fmt.Errorf("%s Exec: %w", operation, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s Affected: %w", operation, err)
	}

	if rowsAffected == 0 {
		return errNoteNotFound
	} else {
		return nil
	}
}

func (n *NoteManagePostgres) UpdateNote(note entities.Note) error {
	const operation = "postgres.UpdateNote"

	stmt, err := n.db.Prepare("UPDATE notes SET title = $1, text = $2, done = $3 WHERE id = $4")
	if err != nil {
		return fmt.Errorf("%s Prepare: %w", operation, err)
	}

	result, err := stmt.Exec(note.Title, note.Text, note.Done, note.Id)
	if err != nil {
		return fmt.Errorf("%s Exec: %w", operation, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s Affected: %w", operation, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s values are not included in the table: %w", operation, err)
	} else {
		return nil
	}
}
