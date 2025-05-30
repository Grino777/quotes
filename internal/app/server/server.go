package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Grino777/quotes/internal/api"
	"github.com/Grino777/quotes/internal/config"
	"github.com/Grino777/quotes/internal/interfaces"
	"github.com/Grino777/quotes/internal/lib/logger"
	serviceAPI "github.com/Grino777/quotes/internal/services/api"
)

const opServer = "app.server."

type ApiProvider interface {
	QuoteProvider
	ApiRouter
}

type ApiRouter interface {
	HomeRoute(w http.ResponseWriter, r *http.Request)
	NotFound(w http.ResponseWriter, r *http.Request)
	NotFoundFallback(w http.ResponseWriter, r *http.Request)
}

type QuoteProvider interface {
	CreateQuote(w http.ResponseWriter, r *http.Request)
	AllQuotes(w http.ResponseWriter, r *http.Request)
	RandomQuote(w http.ResponseWriter, r *http.Request)
	DeleteQuote(w http.ResponseWriter, r *http.Request)
}

type APIServer struct {
	server *http.Server
	logger *slog.Logger
	api    ApiProvider
}

func NewApiServer(log *slog.Logger, cfg *config.APIConfig, storage interfaces.Storage) *APIServer {
	addr := fmt.Sprintf("%s:%s", cfg.Addr, cfg.Port)

	service := serviceAPI.NewService(log, storage)
	apiInstance := api.NewApi(log, service)

	server := &http.Server{Addr: addr}

	return &APIServer{server: server, logger: log, api: apiInstance}
}

func (as *APIServer) Run(ctx context.Context) error {
	const op = opServer + "Run"

	log := as.logger.With("op", op)

	as.setupMultiplexer()
	log.Debug("starting server", slog.String("addr", as.server.Addr))

	errChan := make(chan error, 1)
	go func() {
		if err := as.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("%s: failed to serve: %w", op, err)
		}
	}()

	// Проверка сервера и логирование об успешном запуске если не произошло ошибок
	select {
	case err := <-errChan:
		log.Error("server stopped with error", logger.Error(err))
		return err
	case <-time.After(100 * time.Millisecond):
		log.Info("server started successfully", slog.String("addr", as.server.Addr))
	}

	select {
	case err := <-errChan:
		log.Error("server stopped with error", logger.Error(err))
		return err
	case <-ctx.Done():
		log.Debug("server shutdown initiated due to context cancellation")
		if err := as.Stop(); err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}
		return nil
	}
}

func (as *APIServer) Stop() error {
	const op = opServer + "Stop"

	as.logger.Debug("server shutdown started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := as.server.Shutdown(ctx); err != nil {
		as.logger.Error(fmt.Sprintf("%s: failed to stop server: %v", op, err))
		return err
	}

	as.logger.Debug("server shutdown completed")
	return nil
}

func (as *APIServer) setupMultiplexer() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", as.api.NotFoundFallback)
	mux.HandleFunc("GET /quotes", as.api.AllQuotes)
	mux.HandleFunc("POST /quotes", as.api.CreateQuote)
	mux.HandleFunc("GET /quotes/random", as.api.RandomQuote)
	mux.HandleFunc("DELETE /quotes/", as.api.DeleteQuote)

	middlewares := []func(http.Handler) http.Handler{
		func(h http.Handler) http.Handler { return api.LoggingMiddleware(as.logger, h) },
		api.TimeoutMiddleware(5 * time.Second),
	}

	// Оборачиваем mux в middlewares
	as.server.Handler = api.ApplyMiddlewares(mux, middlewares...)
	as.logger.Debug("all handlers registered")
}
