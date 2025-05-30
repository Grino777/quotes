package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Grino777/quotes/internal/domain/models"
	"github.com/Grino777/quotes/internal/interfaces"
	"github.com/Grino777/quotes/internal/lib/logger"
	"github.com/Grino777/quotes/internal/storage/sqlite"
)

const apiOp = "services.api."

type Service struct {
	logger  *slog.Logger
	storage interfaces.Storage
}

func NewService(log *slog.Logger, storage interfaces.Storage) *Service {
	return &Service{logger: log, storage: storage}
}

func (s *Service) GetQuotes(ctx context.Context) ([]byte, error) {
	const op = apiOp + "GetQuotes"

	log := s.logger.With(slog.String("op", op))

	quotes, err := s.storage.GetQuotes(ctx)
	if err != nil {
		log.Error("failed to get quotes", logger.Error(err))
		return nil, err
	}

	if len(quotes) == 0 {
		return []byte("{}"), nil
	}

	data, err := json.Marshal(quotes)
	if err != nil {
		log.Error("failed to marshaling data", logger.Error(err))
		return nil, err
	}

	return data, nil
}

func (s *Service) CreateQuote(ctx context.Context, quote models.Quote) ([]byte, error) {
	const op = apiOp + "CreateQuote"

	log := s.logger.With(slog.String("op", op))

	res, err := s.storage.CreateQuote(ctx, quote)
	if err != nil {
		if errors.Is(err, sqlite.ErrAlreadyExist) {
			data, err := json.Marshal(map[string]string{"error": "quote already exist"})
			if err != nil {
				log.Error("failed to marshaling data", logger.Error(err))
				return nil, err
			}
			return data, nil
		}
		log.Error("failed to save quote in database", logger.Error(err))
		return nil, err
	}

	data, err := json.Marshal(map[string]any{"result": "success", "id": res})
	if err != nil {
		log.Error("failed to marshaling data", logger.Error(err))
		return nil, err
	}

	return data, nil
}

func (s *Service) GetRandomQuote(ctx context.Context) ([]byte, error) {
	const op = apiOp + "GetRandomQuote"

	log := s.logger.With(slog.String("op", op))

	res, err := s.storage.GetRandomQuote(ctx)
	if err != nil {
		if errors.Is(err, sqlite.ErrQuoteNotExists) {
			return []byte("{}"), nil
		}
		log.Error("failed to get random quote", logger.Error(err))
		return nil, err
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Error("failed to marshaling data", logger.Error(err))
		return nil, err
	}
	return data, nil
}

func (s *Service) FilterQuotes(ctx context.Context, author string) ([]byte, error) {
	const op = apiOp + "CreateQuote"

	log := s.logger.With(slog.String("op", op))

	res, err := s.storage.FilterQuotes(ctx, strings.ToLower(author))
	if err != nil {
		log.Error("failed to get filtered record", logger.Error(err))
		return nil, err
	}

	if len(res) == 0 {
		return []byte("{}"), nil
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Error("failed to marshaling data", logger.Error(err))
		return nil, err
	}

	return data, nil
}

func (s *Service) DeleteQuote(ctx context.Context, id int) ([]byte, error) {
	const op = apiOp + "CreateQuote"

	log := s.logger.With(slog.String("op", op))

	if err := s.storage.DeleteQuote(ctx, id); err != nil {
		if errors.Is(err, sqlite.ErrQuoteNotExists) {
			r := fmt.Sprintf("no quote found with id %d", id)
			data, err := json.Marshal(map[string]string{"error": r})
			if err != nil {
				log.Error("failed to marshaling data", logger.Error(err))
				return nil, err
			}
			return data, nil
		}
		log.Error("failed to get filtered record", logger.Error(err))
		return nil, err
	}

	data, err := json.Marshal(
		map[string]string{"result": fmt.Sprintf("quote with id: %d successfully deleted", id)})
	if err != nil {
		log.Error("failed to marshaling data", logger.Error(err))
		return nil, err
	}
	return data, nil
}
