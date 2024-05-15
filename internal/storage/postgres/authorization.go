package postgres

import (
	"NoteProject/internal/entities"
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
	q, err := s.db.Prepare("INSERT INTO users(name, password) VALUES($1, $2) RETURNING id")

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

		if pgErr.Code.Name() == uniqueViolationCode {
			return 0, errUserExists
		}
		return 0, fmt.Errorf("%s Scan: %w", operation, err)
	}

	return id, nil
}

func (s *AuthPostgres) CheckUser(username, password string) (*entities.User, error) {
	const op = "postgres.CheckUser"

	q, err := s.db.Prepare("SELECT * FROM users WHERE name = $1 AND password = $2")
	if err != nil {
		return nil, fmt.Errorf("%s Prepare: %w", op, err)
	}

	var user entities.User
	rows, err := q.Query(username, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errUserNotExists
		}
		return nil, fmt.Errorf("%s Query: %w", op, err)
	}

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("%s Scan: %w", op, err)
		}
		return &user, nil
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s rows iteration error: %w", op, err)
	}

	return nil, errUserNotExists

}
