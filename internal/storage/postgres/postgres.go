package postgres

import (
	"database/sql"
	"fmt"
)

// PostgresConfig is a struct for managing PG database configuration.
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// Prepare forms a connection string to the PG database based on the configuration.
func (cfg Config) Prepare() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
}

// NewPostgresDB creates a new connection with the PG database.
func NewPostgresDB(cfg Config) (*sql.DB, error) {
	const operation = "storage.NewPostgresDB"

	db, err := sql.Open("postgres", cfg.Prepare())
	if err != nil {
		return nil, fmt.Errorf("%s:%w", operation, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s:%w", operation, err)
	}
	return db, nil
}
