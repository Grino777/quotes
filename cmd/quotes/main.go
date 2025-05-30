package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Grino777/quotes/internal/app"
	"github.com/Grino777/quotes/internal/lib/logger"
)

func main() {
	stop := make(chan os.Signal, 1)
	log := logger.NewLogger(os.Stdout, slog.LevelDebug)

	app, err := app.NewApp(log)
	if err != nil {
		log.Error("failed to create app obj", logger.Error(err))
		return
	}

	if err := app.Run(); err != nil {
		if err := app.Stop(); err != nil {
			log.Error("failed to stop app", logger.Error(err))
		}
		return
	}

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	if err := app.Stop(); err != nil {
		log.Error("failed to stop app", logger.Error(err))
	}
	log.Info("Gracefully stopped")
}
