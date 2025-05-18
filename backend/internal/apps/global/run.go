package global

import (
	"errors"
	"fmt"
	"github.com/goriiin/go-http-balancer/pkg"
	"log/slog"
	"net/http"
)

func (a *App) Run() <-chan error {
	errChan := make(chan error, 1)

	a.router.HandleFunc("/api/v1/data/{id}", a.crud)

	a.router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		pkg.WriteJSON(w, http.StatusOK, struct {
			Answer string `json:"answer"`
		}{"healthy"})

		a.log.Debug("[ health ] health check OK!")
	})

	a.server = &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      a.router,
		ReadTimeout:  a.config.ReadTimeout,
		WriteTimeout: a.config.WriteTimeout,
		IdleTimeout:  a.config.IdleTimeout,
	}

	go func() {
		defer close(errChan)

		a.log.Info("starting server",
			slog.String("address", a.server.Addr),
			slog.Duration("read_timeout", a.config.ReadTimeout),
			slog.Duration("write_timeout", a.config.WriteTimeout),
		)

		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			select {
			case errChan <- fmt.Errorf("server listen error: %w", err):
				a.log.Error("server failed to start", slog.String("error", err.Error()))
			default:
				a.log.Error("failed to send error to channel", slog.String("error", err.Error()))
			}
		}
	}()

	return errChan
}
