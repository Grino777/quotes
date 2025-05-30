package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func ApplyMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func LoggingMiddleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		start := time.Now()
		next.ServeHTTP(w, r)
		execTime := time.Since(start).Seconds()
		log.Debug("request executed",
			slog.String("path", r.URL.Path),
			slog.Float64("exec_time_sec", execTime),
		)
	})
}

func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			r = r.WithContext(ctx)
			done := make(chan struct{})
			var respErr error
			go func() {
				defer close(done)
				next.ServeHTTP(w, r)
			}()
			select {
			case <-done:
				return
			case <-ctx.Done():
				respErr = ctx.Err()
				if respErr == context.DeadlineExceeded {
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}
		})
	}
}
