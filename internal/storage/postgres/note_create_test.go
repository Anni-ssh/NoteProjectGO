package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateNote(T *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		T.Fatal(err)
	}
	defer db.Close()

	s := NewNoteManage(db)

	type note struct {
		userID int
		title  string
		text   string
	}

	type mockBehaviour func(n note, id int)

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		note          note
		result        int
	}{
		{name: "Success",
			note: note{userID: 70,
				title: "Заметки",
				text:  "Мои заметки"},
			result: 5,
			mockBehaviour: func(n note, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO notes").WithArgs(n.userID, n.title, n.text).WillReturnRows(rows)
			},
		}}

	for _, tt := range testTable {
		T.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.note, tt.result)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := s.CreateNote(ctx, tt.note.userID, tt.note.title, tt.note.text)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.result, result)
			}
		})
	}
}
