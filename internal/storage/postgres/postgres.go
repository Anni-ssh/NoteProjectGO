package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func (cfg Config) Prepare() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
}

func NewPostgresDB(cfg Config, attempts int, delay time.Duration) (*sql.DB, error) {
	const operation = "storage.NewPostgresDB"

	var db *sql.DB
	var err error

	for i := 0; i < attempts; i++ {
		db, err = sql.Open("postgres", cfg.Prepare())
		if err != nil {
			if i < attempts-1 {
				time.Sleep(delay)
				continue
			}
			return nil, fmt.Errorf("%s - failed to open database connection: %w", operation, err)
		}

		if err = db.Ping(); err != nil {
			if i < attempts-1 {
				time.Sleep(delay)
				continue
			}
			return nil, fmt.Errorf("%s - failed to ping database: %w", operation, err)
		}

		return db, nil
	}

	// Не должно быть достигнуто
	return nil, fmt.Errorf("%s - failed to connect to database after %d attempts", operation, attempts)
}
