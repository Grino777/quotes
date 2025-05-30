package sqlite

import (
	"database/sql"
	"log/slog"

	"github.com/Grino777/quotes/internal/config"
	"github.com/Grino777/quotes/internal/lib/logger"
	_ "github.com/mattn/go-sqlite3"
)

const sqliteChema = `
	CREATE TABLE IF NOT EXISTS quotes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	author VARCHAR(100) NOT NULL,
	quote TEXT NOT NULL,
	CONSTRAINT unique_quote UNIQUE (author, quote)
	);
`

const sqliteOp = "storage.sqlite."

type Storage struct {
	logger *slog.Logger
	cfg    *config.SQLiteConfig
	client *sql.DB
}

func NewStorage(
	log *slog.Logger,
	cfg *config.SQLiteConfig,
) *Storage {
	return &Storage{
		logger: log,
		cfg:    cfg,
	}
}

func (s *Storage) Connect() error {
	const op = sqliteOp + "Connect"

	log := s.logger.With(slog.String("op", op))

	conn, err := sql.Open("sqlite3", s.cfg.Addr)
	if err != nil {
		log.Error("failed to connect database", logger.Error(err))
		return err
	}

	if _, err := conn.Exec(sqliteChema); err != nil {
		log.Error("failed to create table for database", logger.Error(err))
		return err
	}

	s.client = conn
	return nil
}

func (s *Storage) Close() error {
	const op = sqliteOp + "Close"

	if s.client != nil {
		if err := s.client.Close(); err != nil {
			s.logger.Error("failed to close database connection", slog.String("op", op), logger.Error(err))
			return err
		}
	}
	return nil
}
