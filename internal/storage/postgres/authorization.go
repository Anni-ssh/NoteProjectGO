package postgres

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/errs"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (s *AuthPostgres) CreateUser(user entities.User) (int, error) {
	const operation = "postgres.CreateUser"
	q, err := s.db.Prepare("INSERT INTO users(name, password_hash) VALUES($1, $2) RETURNING id")

	if err != nil {
		return 0, fmt.Errorf("%s Prepare: %w", operation, err)
	}

	var id int

	err = q.QueryRow(user.Username, user.Password).Scan(&id)

	if err != nil {
		var pgErr *pq.Error
		ok := errors.As(err, &pgErr)
		if !ok {
			return 0, fmt.Errorf("%s Scan: %w", operation, err)
		}
		// Проверка ошибки уникальность значения
		if pgErr.Code.Name() == "unique_violation" {
			return 0, errs.ErrUserExists
		}
		return 0, fmt.Errorf("%s Scan: %w", operation, err)
	}

	return id, nil
}

func (s *AuthPostgres) CheckUser(username, password string) (entities.User, error) {
	const op = "postgres.CheckUser"

	var user entities.User

	q, err := s.db.Prepare("SELECT * FROM users WHERE name = $1 AND password_hash = $2")
	if err != nil {
		return user, fmt.Errorf("%s Prepare: %w", op, err)
	}

	rows, err := q.Query(username, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errs.ErrUserNotExists
		}
		return user, fmt.Errorf("%s Query: %w", op, err)
	}

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			return user, fmt.Errorf("%s Scan: %w", op, err)
		}
	}

	if err = rows.Err(); err != nil {
		return user, fmt.Errorf("%s rows iteration error: %w", op, err)
	}

	return user, nil

}
