// Здесь описываются запросы в БД

package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	// "github.com/shameoff/more-than-trip/core/internal/domain/models"
	// "github.com/shameoff/more-than-trip/core/internal/jaeger"
	// "github.com/shameoff/more-than-trip/core/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}
