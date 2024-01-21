package helperFunc

import (
	"TestProject/internal/lib/dataBaseSQL"
	"TestProject/internal/lib/dataBaseSQL/ServSQLite"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrDataNil = errors.New("data not found")
)

// FIX ME name ctx
// СonvExtractedNotesData подготавливает заметки к выводу на страницу.
func СonvExtractedNotesData(DB *sql.DB, ctx context.Context) (*dataBaseSQL.NotesList, error) {

	const operation = "helperFunc.СonvExtractedNotesData" //FIX ME name
	typeNote := ServSQLite.DataBaseSQLiteNote{Storage: DB}

	bodyNote, err := typeNote.CheckNotesList(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", operation, err)
	}

	return bodyNote, nil

}

// FIX ME name ctx
// Возможно в параметр стот добавить список
// SendNoteToDB подготавливает заметки к выводу на страницу.
func SendNoteToDB(DB *sql.DB, ctx context.Context, notesList ...dataBaseSQL.Note) error {

	const operation = "helperFunc.SendNoteToDB" //FIX ME name
	typeNote := ServSQLite.DataBaseSQLiteNote{Storage: DB}

	//FIX ME транзакции SQL
	for _, val := range notesList {
		err := typeNote.SaveNoteInDB(ctx, val)
		if err != nil {
			return fmt.Errorf("%s:%w", operation, err)
		}
	}
	return nil
}

func ConvStrInNote(title, text string) dataBaseSQL.Note {
	newNote := dataBaseSQL.Note{
		Title: title,
		Text:  text,
	}
	return newNote
}
