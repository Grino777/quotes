package app

import (
	"context"
	"log/slog"
	"sync"

	"github.com/Grino777/quotes/internal/app/server"
	"github.com/Grino777/quotes/internal/config"
	"github.com/Grino777/quotes/internal/interfaces"
	"github.com/Grino777/quotes/internal/lib/logger"
	"github.com/Grino777/quotes/internal/storage/sqlite"
	sqliteU "github.com/Grino777/quotes/internal/utils/sqlite"
)

const opApp = "app."

type App struct {
	Logger    *slog.Logger
	Config    *config.Config
	ApiServer *server.APIServer
	Storage   interfaces.Storage
	cancel    context.CancelFunc
}

func NewApp(log *slog.Logger) (*App, error) {
	const op = opApp + "NewApp"

	config, err := config.NewConfig()
	if err != nil {
		log.Error("failed to get configs for app", slog.String("op", op), logger.Error(err))
		return nil, err
	}

	storage := sqlite.NewStorage(log, &config.SQLite)
	server := server.NewApiServer(log, &config.API, storage)

	return &App{
		Logger:    log,
		Config:    config,
		ApiServer: server,
		Storage:   storage,
	}, nil
}

func (a *App) Run() error {
	const op = opApp + "Run"

	log := a.Logger.With("op", op)
	var (
		wg      sync.WaitGroup
		lastErr error
	)
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	if err := sqliteU.CheckStorageFolder(); err != nil {
		log.Error("failed to check storage dir", logger.Error(err))
		return err
	}

	if err := a.Storage.Connect(); err != nil {
		log.Error("failed to connect to database", logger.Error(err))
		return err
	}

	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := a.ApiServer.Run(ctx); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Debug("app shutdown initiated")
	case err := <-errChan:
		log.Error("stopping app due to error", logger.Error(err))
		lastErr = err
	}

	wg.Wait()

	return lastErr
}

func (a *App) Stop() error {
	a.cancel()

	if a.ApiServer != nil {
		if err := a.ApiServer.Stop(); err != nil {
			return err
		}
	}

	if a.Storage != nil {
		if err := a.Storage.Close(); err != nil {
			return err
		}
	}

	a.Logger.Debug("app successfully stopped")
	return nil
}
