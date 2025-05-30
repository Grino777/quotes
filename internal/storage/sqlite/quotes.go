package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Grino777/quotes/internal/domain/models"
	"github.com/mattn/go-sqlite3"
)

var (
	ErrAlreadyExist   = errors.New("quote already exists")
	ErrQuoteNotExists = errors.New("quote not exists")
)

const (
	ReqDuration = 5
)

const opQuotes = "storage.sqlite."

func (s *Storage) GetQuotes(ctx context.Context) ([]models.Quote, error) {
	const op = opQuotes + "GetQuotes"

	ctx, cancel := context.WithTimeout(ctx, ReqDuration*time.Second)
	defer cancel()

	stmt := `SELECT id, author, quote FROM quotes`

	rows, err := s.client.QueryContext(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to getting all quotes: %w", op, err)
	}
	defer rows.Close()

	var quotes []models.Quote

	for rows.Next() {
		var q models.Quote
		if err := rows.Scan(&q.Id, &q.Author, &q.Quote); err != nil {
			return nil, fmt.Errorf("%s: failed to scanning result: %w", op, err)
		}
		quotes = append(quotes, q)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: failed to processing query result: %w", op, err)
	}

	return quotes, nil
}

func (s *Storage) CreateQuote(ctx context.Context, quote models.Quote) (int64, error) {
	const op = opQuotes + "CreateQuote"

	ctx, cancel := context.WithTimeout(ctx, ReqDuration*time.Second)
	defer cancel()

	stmt := `INSERT INTO quotes (author, quote) VALUES (?, ?)`

	result, err := s.client.ExecContext(ctx, stmt, strings.ToLower(quote.Author), quote.Quote)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			return 0, ErrAlreadyExist
		}
		return 0, fmt.Errorf("%s: failed to insert quote: %w", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to retrieve last insert ID: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetRandomQuote(ctx context.Context) (models.Quote, error) {
	const op = opQuotes + "GetRandomQuote"

	ctx, cancel := context.WithTimeout(ctx, ReqDuration*time.Second)
	defer cancel()

	stmt := `SELECT id, author, quote FROM quotes ORDER BY RANDOM() LIMIT 1`

	row := s.client.QueryRowContext(ctx, stmt)

	var q models.Quote
	if err := row.Scan(&q.Id, &q.Author, &q.Quote); err != nil {
		if err == sql.ErrNoRows {
			return models.Quote{}, ErrQuoteNotExists
		}
		return models.Quote{}, fmt.Errorf("%s: failed to scan random quote: %w", op, err)
	}

	return q, nil
}

func (s *Storage) DeleteQuote(ctx context.Context, id int) error {
	const op = opQuotes + "DeleteQuote"

	ctx, cancel := context.WithTimeout(ctx, ReqDuration*time.Second)
	defer cancel()

	stmt := `DELETE FROM quotes WHERE id = ?`

	result, err := s.client.ExecContext(ctx, stmt, id)
	if err != nil {
		return fmt.Errorf("%s: failed to delete quote: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to retrieve rows affected: %w", op, err)
	}
	if rowsAffected == 0 {
		return ErrQuoteNotExists
	}

	return nil
}

func (s *Storage) FilterQuotes(ctx context.Context, author string) ([]models.Quote, error) {
	const op = opQuotes + "FilterQuotes"

	ctx, cancel := context.WithTimeout(ctx, ReqDuration*time.Second)
	defer cancel()

	stmt := `SELECT id, author, quote FROM quotes WHERE author = ?`

	rows, err := s.client.QueryContext(ctx, stmt, author)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query quotes by author: %w", op, err)
	}
	defer rows.Close()

	var quotes []models.Quote

	for rows.Next() {
		var q models.Quote
		if err := rows.Scan(&q.Id, &q.Author, &q.Quote); err != nil {
			return nil, fmt.Errorf("%s: failed to scan quote: %w", op, err)
		}
		quotes = append(quotes, q)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: failed to process query result: %w", op, err)
	}

	return quotes, nil
}
