package postgres

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/errs"
	"context"
	"database/sql"
	"fmt"
)

type NoteManagePostgres struct {
	db *sql.DB
}

func NewNoteManage(db *sql.DB) *NoteManagePostgres {
	return &NoteManagePostgres{db: db}
}

// CreateNote создает новую заметку для пользоватля по его ID.
func (n *NoteManagePostgres) CreateNote(ctx context.Context, userID int, title, text string) (int, error) {
	const op = "postgres.CreateNote"

	// Начало транзакции
	tx, err := n.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: BeginTx failed: %w", op, err)
	}

	// Проверка существования пользователя
	var userExists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&userExists)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: Check user exists failed: %w", op, err)
	}
	if !userExists {
		tx.Rollback()
		return 0, fmt.Errorf("%s: User does not exist: %w", op, errs.ErrUserNotExists)
	}

	// Подготовка SQL-запроса
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO notes (id, title, text) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: Prepare failed: %w", op, err)
	}

	var id int
	// Добавление заметки
	err = stmt.QueryRowContext(ctx, userID, title, text).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: Insert failed: %w", op, err)
	}

	// Коммит транзакции
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: Commit failed: %w", op, err)
	}

	return id, nil
}

// NotesList возвращает список заметок по ID пользователя.
func (n *NoteManagePostgres) NotesList(ctx context.Context, userID int) ([]entities.Note, error) {
	const op = "postgres.NotesList"

	// Начало транзакции
	tx, err := n.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: BeginTx failed: %w", op, err)
	}

	// Проверка наличия пользователя
	var userExists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&userExists)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: Check user exists failed: %w", op, err)
	}
	if !userExists {
		tx.Rollback()
		return nil, fmt.Errorf("%s: User does not exist: %w", op, errs.ErrUserNotExists)
	}

	// Подготовка SQL-запроса
	q, err := tx.PrepareContext(ctx, "SELECT * FROM notes WHERE id = $1")
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: Prepare failed: %w", op, err)
	}

	// Выполнение запроса
	rows, err := q.QueryContext(ctx, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: Query failed: %w", op, err)
	}

	// Чтение результатов
	var notes []entities.Note
	for rows.Next() {
		var note entities.Note
		err := rows.Scan(&note.Id, &note.UserId, &note.Title, &note.Text, &note.Date, &note.Done)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%s: Scan failed: %w", op, err)
		}
		notes = append(notes, note)
	}

	// Проверка наличия ошибок
	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: Rows iteration error: %w", op, err)
	}

	// Коммит
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: Commit failed: %w", op, err)
	}

	return notes, nil
}

// UpdateNote обновляет заметку по ID, который указан в entities.Note
func (n *NoteManagePostgres) UpdateNote(ctx context.Context, note entities.Note) error {
	const operation = "postgres.UpdateNote"

	// Подготовка SQL-запроса
	stmt, err := n.db.PrepareContext(ctx, "UPDATE notes SET title = $1, text = $2, done = $3 WHERE id = $4")
	if err != nil {
		return fmt.Errorf("%s: Prepare failed: %w", operation, err)
	}

	// Выполнение SQL-запроса
	// Не выполняется проверка наличия заметки т.к. ExecContext вернёт кол-во затронутых row.
	result, err := stmt.ExecContext(ctx, note.Title, note.Text, note.Done, note.Id)
	if err != nil {
		return fmt.Errorf("%s: Exec failed: %w", operation, err)
	}

	// Получение количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: RowsAffected failed: %w", operation, err)
	}

	// Если количество затронутых строк 0, то заметка не найдена
	if rowsAffected == 0 {
		return fmt.Errorf("%s: No rows were updated: %w", operation, errs.ErrNoteNotExists)
	}

	return nil
}

// DeleteNote удаляет заметку по её ID.
func (n *NoteManagePostgres) DeleteNote(ctx context.Context, noteID int) error {
	const operation = "postgres.DeleteNote"

	// Подготовка SQL-запроса
	q, err := n.db.PrepareContext(ctx, "DELETE FROM notes WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: PrepareContext failed: %w", operation, err)
	}

	// Выполнение SQL-запроса
	result, err := q.ExecContext(ctx, noteID)
	if err != nil {
		return fmt.Errorf("%s: ExecContext failed: %w", operation, err)
	}

	// Получение затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: RowsAffected failed: %w", operation, err)
	}

	// Если количество затронутых строк 0, то заметка не найдена
	if rowsAffected == 0 {
		return errs.ErrNoteNotExists
	}

	return nil
}
