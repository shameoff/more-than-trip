package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HttpApp struct {
	log    *slog.Logger
	server *http.Server
	port   int
	Router *chi.Mux
}

// New создает новый экземпляр HTTP-сервера с заданными настройками
func New(log *slog.Logger, port int) *HttpApp {
	router := chi.NewRouter()

	// Добавление middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Пример обработчика для health-check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return &HttpApp{
		log: log,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		},
		port:   port,
		Router: router,
	}
}

// MustRun запускает HTTP-сервер и паникует в случае ошибки
func (a *HttpApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run запускает HTTP-сервер
func (a *HttpApp) Run() error {
	a.log.Info("HTTP server is running", slog.String("address", a.server.Addr))

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("httpapp.Run: %w", err)
	}
	return nil
}

// Stop останавливает HTTP-сервер с использованием GracefulShutdown
func (a *HttpApp) Stop(ctx context.Context) error {
	a.log.Info("stopping HTTP server", slog.Int("port", a.port))
	return a.server.Shutdown(ctx)
}
