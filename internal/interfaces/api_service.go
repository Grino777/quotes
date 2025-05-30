package interfaces

import (
	"context"

	"github.com/Grino777/quotes/internal/domain/models"
)

type Service interface {
	GetQuotes(ctx context.Context) ([]byte, error)
	CreateQuote(ctx context.Context, quote models.Quote) ([]byte, error)
	GetRandomQuote(ctx context.Context) ([]byte, error)
	FilterQuotes(ctx context.Context, author string) ([]byte, error)
	DeleteQuote(ctx context.Context, id int) ([]byte, error)
}
