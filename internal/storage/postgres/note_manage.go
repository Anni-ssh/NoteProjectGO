package postgres

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/errs"
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
	const op = "postgres.CreateNote"

	// Начало транзакции
	tx, err := n.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	// Проверка существования пользователя
	var userExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)", userID).Scan(&userExists)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s check user exists: %w", op, err)
	}
	if !userExists {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, errs.ErrUserNotExists)
	}

	stmt, err := tx.Prepare("INSERT INTO notes (user_id, title, text) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s Prepare: %w", op, err)
	}

	var id int
	err = stmt.QueryRow(userID, title, text).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s Insert: %w", op, err)
	}

	// Коммит транзакции
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s Commit: %w", op, err)
	}

	return id, nil
}

func (n *NoteManagePostgres) NotesList(userID int) ([]entities.Note, error) {
	const op = "postgres.NotesList"

	// Проверка наличия пользователя
	// Позволяет корректно отличить отсутстивие заметок и отсутствие пользователя
	// Начало транзакции
	tx, err := n.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("database error: %w, operation: %s", err, op)
	}

	// Провека наличия пользователя
	var userExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)", userID).Scan(&userExists)
	if err != nil {
		return nil, fmt.Errorf("%s check user exists: %w", op, err)
	}
	if !userExists {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrUserNotExists)
	}

	// Подготвка и исполнение запроса
	q, err := tx.Prepare("SELECT * FROM notes WHERE user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s Prepare: %w", op, err)
	}

	rows, err := q.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("%s Query: %w", op, err)
	}

	var notes []entities.Note
	for rows.Next() {
		var note entities.Note
		err := rows.Scan(&note.Id, &note.UserId, &note.Title, &note.Text, &note.Date, &note.Done)
		if err != nil {
			return nil, fmt.Errorf("%s Scan: %w", op, err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s Rows: %w", op, err)
	}

	// Коммит
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s Commit: %w", op, err)
	}

	return notes, nil
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
		return fmt.Errorf("%s RowsAffected: %w", operation, err)
	}

	// Если кол-во записей 0, то такой заметки нет
	if rowsAffected == 0 {
		return fmt.Errorf("%s no rows were updated: %w", operation, errs.ErrNoteNotExists)
	}

	return nil
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
		return errs.ErrNoteNotExists
	} else {
		return nil
	}
}
