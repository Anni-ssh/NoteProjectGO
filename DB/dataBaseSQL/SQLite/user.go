package SQLite

import (
	"Test/dataBaseSQL"
	"Test/dataBaseSQL/errSql"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

type DataBaseSQLiteUser struct {
	Storage *sql.DB
}

// OpenDB Методы
// OpenDB открытие соединения с базой данных, возвращает открытую БД.
func (dataBase *DataBaseSQLiteUser) OpenDB(path string) error {
	const operation = "DataBaseSQLiteUser.OpenDB"
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

// CreateUsersTable создаёт новую таблицу users, если она не существует.
func (dataBase DataBaseSQLiteUser) CreateUsersTable(ctx context.Context) error {
	const operation = "DataBaseSQLiteUser.CreateUsersTable"

	UsersQuery := `CREATE TABLE IF NOT EXISTS users
    (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL,
        age INTEGER NOT NULL,
        password TEXT NOT NULL,
        superUser INTEGER DEFAULT 0
    );`

	req, err := dataBase.Storage.PrepareContext(ctx, UsersQuery)
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

// SaveUserInDB добавление значений users в базу данных.
func (dataBase DataBaseSQLiteUser) SaveUserInDB(ctx context.Context, user dataBaseSQL.User) error {
	const operation = "DataBaseSQLiteUser.SaveUserInDB"

	query, args, err := squirrel.Insert("users").
		Columns("name", "age", "password").
		Values(user.Name, user.Age, user.Password).
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

// CheckUser запрос для проверки и извлечения полей для структуры User из БД, если нет такого пользователя, то возвращается nil
func (dataBase DataBaseSQLiteUser) CheckUser(ctx context.Context, name, password string) (*dataBaseSQL.User, error) {
	const operation = "DataBaseSQLiteUser.CheckUser"

	query, args, err := squirrel.Select("id", "name", "age", "password", "superuser").
		From("users").
		Where(squirrel.Eq{"name": name, "password": password}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: ошибка при создании sql запроса: %w", operation, err)
	}

	req, err := dataBase.Storage.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка при подготовке запроса: %w", operation, err)
	}
	var user dataBaseSQL.User

	err = req.QueryRowContext(ctx, args...).Scan(&user.Id, &user.Name, &user.Age, &user.Password, &user.SuperUser)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: пользователь не найден: %w", operation, errSql.ErrUserNotFound)
		}

		if ctx.Err() != nil {
			return nil, fmt.Errorf("%s: превышено время ожидания ответа: %w", operation, ctx.Err())
		}

		return nil, fmt.Errorf("%s: ошибка при выполнении запроса: %w", operation, err)

	}

	return &user, nil
}

// DeleteUser удаляет определённого пользователя
func (dataBase DataBaseSQLiteUser) DeleteUser(ctx context.Context, id int) error {
	const operation = "DataBaseSQLiteUser.DeleteUser"

	query, args, err := squirrel.Delete("users").Where(squirrel.Eq{"id": id}).ToSql()

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
