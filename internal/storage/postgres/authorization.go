package postgres

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/errs"
	"context"
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

// CreateUser создает нового пользователя.
func (s *AuthPostgres) CreateUser(ctx context.Context, user entities.User) (int, error) {
	const operation = "postgres.CreateUser"

	// Подготовка SQL-запроса
	q, err := s.db.PrepareContext(ctx, "INSERT INTO users(name, password_hash) VALUES($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: Prepare failed: %w", operation, err)
	}

	var id int

	// Выполнение запроса
	err = q.QueryRowContext(ctx, user.Username, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			// Проверка на ошибку уникальности значения
			if pgErr.Code.Name() == "unique_violation" {
				return 0, fmt.Errorf("%s: user with username %s already exists: %w", operation, user.Username, errs.ErrUserExists)
			}
			return 0, fmt.Errorf("%s: PostgreSQL error: %s", operation, pgErr.Message)
		}
		return 0, fmt.Errorf("%s: Scan failed: %w", operation, err)
	}

	return id, nil
}

// CheckUser проверяет существование пользователя.
func (s *AuthPostgres) CheckUser(ctx context.Context, username, password string) (entities.User, error) {
	const op = "postgres.CheckUser"

	var user entities.User

	// Подготовка SQL-запроса
	q, err := s.db.PrepareContext(ctx, "SELECT * FROM users WHERE name = $1 AND password_hash = $2")
	if err != nil {
		return user, fmt.Errorf("%s: Prepare failed: %w", op, err)
	}

	// Выполнение запроса
	row := q.QueryRowContext(ctx, username, password)

	// Чтение результатов
	err = row.Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, errs.ErrUserNotExists)
		}
		return user, fmt.Errorf("%s: Scan failed: %w", op, err)
	}

	return user, nil

}
