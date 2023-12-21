package dataBaseSQL

import (
	"context"
	"database/sql"
	"fmt"
)

type DataBase interface {
	Users
	Notes
}

type Users interface {
	CreateUsersTable(ctx context.Context) error
	SaveUserInDB(ctx context.Context, user User) error
	CheckUser(ctx context.Context, name, password string) (*User, error)
	DeleteUser(ctx context.Context, id int) error
}

type Notes interface {
	CreateNotesTable(ctx context.Context) error
	SaveNoteInDB(ctx context.Context, note Note) error
	CheckNoteByID(ctx context.Context, id int) (*Note, error)
	DeleteNote(ctx context.Context, id int) error
}

func OpenDB(path string) (*sql.DB, error) {
	const operation = "FuncOpenDB"
	req, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", operation, err)
	}
	if err = req.Ping(); err != nil {
		return nil, fmt.Errorf("%s:%w", operation, err)
	}
	return req, nil
}

// Создание объекта users
type User struct {
	Id        int    `json:"id" env-required:"false"`
	Name      string `json:"name" env-required:"true"`
	Age       int    `json:"age" env-required:"true"`
	Password  string `json:"password" env-required:"true"`
	SuperUser bool   `json:"superUser" env-required:"false"`
}

// МЕТОДЫ USER
// UserInfo структурированный вывод данных о пользователе.
func (u User) UserInfo() string {
	return fmt.Sprintf("Id:%d\nName:%s\nAge:%d\nPassword:%s\nSuperUser:%t\n", u.Id, u.Name, u.Age, u.Password, u.SuperUser)
}

// Создание объекта note
type Note struct {
	Id     int    `json:"id" env-required:"true"`
	IdUser int    `json:"id_user" env-required:"true"`
	Title  string `json:"title" env-required:"true"`
	Text   string `json:"text" env-required:"true"`
	Done   bool   `json:"done" env-required:"true"`
}

// МЕТОДЫ NOTE
// NoteInfo структурированный вывод данных о заметке.
func (n Note) NoteInfo() string {
	return fmt.Sprintf("Id:%d\nIdUser:%d\nTitle:%s\nText:%s\nDone:%t\n", n.Id, n.IdUser, n.Title, n.Text, n.Done)
}
