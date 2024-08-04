package handler

import (
	service2 "NoteProject/internal/service"
	mock_service "NoteProject/internal/service/mocks"
	"NoteProject/pkg/logger"
	"context"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_create(t *testing.T) {
	type mockBehavior func(s *mock_service.MockNote, note inputNote)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testTable := []struct {
		name                 string
		inputBody            string
		inputNote            inputNote
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{name: "OK",
			inputBody: `{"user_id": 1, "title": "Test", "text": "Test save"}`,
			inputNote: inputNote{
				UserId: 1,
				Title:  "Test",
				Text:   "Test save",
			},
			mockBehavior: func(s *mock_service.MockNote, note inputNote) {
				s.EXPECT().CreateNote(ctx, note.UserId, note.Title, note.Text).Return(1, nil)

			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1}`},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noteManage := mock_service.NewMockNote(c)
			tt.mockBehavior(noteManage, tt.inputNote)

			service := &service2.Service{
				Note: noteManage,
			}

			handler := Handler{services: service, Logs: logger.SetupLogger("local", os.Stdout)}

			r := chi.NewRouter()
			r.Post("/note/create", handler.noteCreate)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/note/create", strings.NewReader(tt.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.TrimSpace(w.Body.String()))

		})
	}

}
