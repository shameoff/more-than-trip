package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/shameoff/more-than-trip/core/internal/app"
	"github.com/shameoff/more-than-trip/core/internal/config"
	"github.com/shameoff/more-than-trip/core/internal/lib/logger/sl"
)

func main() {
	// Инициализируем объект конфига
	cfg := config.MustLoad()

	// Инициализируем логгер
	log := sl.New(cfg.Env)

	// Инициализируем приложение (app)
	application := app.New(log, cfg)
	ctx := context.Background()
	go func() {
		application.HTTPServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	// Waiting for SIGINT (pkill -2) or SIGTERM
	<-stop

	// initiate graceful shutdown
	application.HTTPServer.Stop(ctx) // Assuming GRPCServer has Stop() method for graceful shutdown
	log.Info("Gracefully stopped")
	os.Exit(0)
}
