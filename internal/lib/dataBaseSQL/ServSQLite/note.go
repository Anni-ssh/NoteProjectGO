package ServSQLite

import (
	"TestProject/internal/lib/dataBaseSQL"
	"TestProject/internal/lib/dataBaseSQL/errSql"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type DataBaseSQLiteNote struct {
	Storage *sql.DB
}

// Методы
// OpenDB открытие соединения с базой данных, возвращает открытую БД.
func (dataBase *DataBaseSQLiteNote) OpenDB(path string) error {
	const operation = "DataBaseSQLiteNote.OpenDB"
	req, err := sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("%s:%w", operation, err)
	}
	if err = req.Ping(); err != nil {
		return fmt.Errorf("%s:%w", operation, err)
	}
	dataBase.Storage = req
	return nil
}

// CreateTables создаёт новую таблицу notes, если она не существует.
func (dataBase DataBaseSQLiteNote) CreateNotesTable(ctx context.Context) error {
	const operation = "DataBaseSQLiteNote.CreateNotesTable"

	NotesQuery := `CREATE TABLE IF NOT EXISTS notes
    (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER DEFAULT 0,
        title TEXT NOT NULL,
        text TEXT,
        date TEXT,
        done INTEGER DEFAULT 0
    );`

	req, err := dataBase.Storage.PrepareContext(ctx, NotesQuery)
	if err != nil {
		return fmt.Errorf("%s: ошибка при подготовке запроса: %w", operation, err)
	}

	_, err = req.ExecContext(ctx)

	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("%s: превышено время ожидания ответа: %w", operation, ctx.Err())
		}
		return fmt.Errorf("%s: ошибка при выполнении запроса: %w", operation, err)
	}

	return nil
}

// SaveNoteInDB добавление значений notes в базу данных.
func (dataBase DataBaseSQLiteNote) SaveNoteInDB(ctx context.Context, note dataBaseSQL.Note) error {
	const operation = "DataBaseSQLiteNote.SaveNoteInDB"
	//Фиксация текущего времени
	currentDate := time.Now().Format("2006-01-02 15:04:05")

	query, args, err := squirrel.Insert("notes").
		Columns("user_id", "title", "text", "date", "done").
		Values(note.IdUser, note.Title, note.Text, currentDate, note.Done).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: ошибка при создании sql запроса: %w", operation, err)
	}

	req, err := dataBase.Storage.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: ошибка при подготовке запроса: %w", operation, err)
	}

	_, err = req.ExecContext(ctx, args...)

	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("%s: превышено время ожидания ответа: %w", operation, ctx.Err())
		}
		return fmt.Errorf("%s: ошибка при выполнении запроса: %w", operation, err)
	}
	return nil
}

// CheckNoteByID запрос для проверки и извлечения полей для структуры Note из БД.
func (dataBase DataBaseSQLiteNote) CheckNoteByID(ctx context.Context, id int) (*dataBaseSQL.Note, error) {
	const operation = "DataBaseSQLiteNote.CheckNoteByID"

	query, args, err := squirrel.Select("id", "user_id", "title", "text", "date", "done").
		From("notes").
		Where(squirrel.Eq{"id": id}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: ошибка при создании sql запроса: %w", operation, err)
	}

	req, err := dataBase.Storage.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка при подготовке запроса: %w", operation, err)
	}
	var note dataBaseSQL.Note

	err = req.QueryRowContext(ctx, args...).Scan(&note.Id, &note.IdUser, &note.Title, &note.Text, &note.Date, &note.Done)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: заметка не найдена: %w", operation, errSql.ErrNoteNotFound)
		}

		if ctx.Err() != nil {
			return nil, fmt.Errorf("%s: превышено время ожидания ответа: %w", operation, ctx.Err())
		}

		return nil, fmt.Errorf("%s: ошибка при выполнении запроса: %w", operation, err)

	}

	return &note, nil
}

// CheckNotesList вывод всех значений из БД и создание списка NotesList.
func (dataBase DataBaseSQLiteNote) CheckNotesList(ctx context.Context) (*dataBaseSQL.NotesList, error) {
	const operation = "DataBaseSQLiteNote.CheckNotesList"

	query, _, err := squirrel.Select("id", "user_id", "title", "text", "date", "done").
		From("notes").ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: ошибка при создании sql запроса: %w", operation, err)
	}

	req, err := dataBase.Storage.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка при подготовке запроса: %w", operation, err)
	}

	//FIX ME нужно ли закрывать данные соединения
	defer req.Close()

	rows, err := req.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка при выполнении запроса: %w", operation, err)
	}
	//FIX ME нужно ли закрывать данные соединения
	defer rows.Close()

	var notesList dataBaseSQL.NotesList

	for rows.Next() {
		var note dataBaseSQL.Note
		err = rows.Scan(&note.Id, &note.IdUser, &note.Title, &note.Text, &note.Date, &note.Done)

		if err != nil {
			return nil, fmt.Errorf("%s: ошибка при сканировании строки результата: %w", operation, err)
		}

		notesList.Notes = append(notesList.Notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: ошибка после завершения чтения строк: %w", operation, err)
	}

	return &notesList, nil
}

// DeleteNote удаляет определённую заметку по id
func (dataBase DataBaseSQLiteNote) DeleteNote(ctx context.Context, id int) error {
	const operation = "DataBaseSQLiteNote.DeleteNote"

	query, args, err := squirrel.Delete("notes").Where(squirrel.Eq{"id": id}).ToSql()

	if err != nil {
		return fmt.Errorf("%s: ошибка при создании sql запроса: %w", operation, err)
	}

	req, err := dataBase.Storage.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: ошибка при подготовке sql запроса: %w", operation, err)
	}

	_, err = req.ExecContext(ctx, args...)
	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("%s: превышено время ожидания ответа: %w", operation, err)
		}
		return fmt.Errorf("%s: Ошибка удаления пользователя из базы данных: %w", operation, err)
	}
	return nil
}
