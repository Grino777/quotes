package interfaces

import (
	"context"

	"github.com/Grino777/quotes/internal/domain/models"
)

type Storage interface {
	GetQuotes(ctx context.Context) ([]models.Quote, error)
	CreateQuote(ctx context.Context, quote models.Quote) (int64, error)
	GetRandomQuote(ctx context.Context) (models.Quote, error)
	FilterQuotes(ctx context.Context, author string) ([]models.Quote, error)
	DeleteQuote(ctx context.Context, id int) error
	Connect() error
	Close() error
}
