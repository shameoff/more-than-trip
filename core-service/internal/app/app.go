package app

import (
	"log/slog"

	httpapp "github.com/shameoff/more-than-trip/core/internal/app/http"
	"github.com/shameoff/more-than-trip/core/internal/config"
	"github.com/shameoff/more-than-trip/core/internal/http-server/core"
	routes "github.com/shameoff/more-than-trip/core/internal/http-server/core"
	"github.com/shameoff/more-than-trip/core/internal/lib/converter"
	coreService "github.com/shameoff/more-than-trip/core/internal/services/core"
	"github.com/shameoff/more-than-trip/core/internal/storage/postgres"
)

type App struct {
	HTTPServer *httpapp.HttpApp
}

func New(log *slog.Logger,
	config *config.Config) *App {
	const op = "app.New"

	log = log.With(
		slog.String("op", op),
	)

	// Init Database connection
	pgDSN, err := converter.ConvertDatabaseConfigToDSN(config.Database)
	if err != nil {
		panic(err)
	}
	storage, err := postgres.New(pgDSN)
	log.Debug("successfully connected to database")

	// Init core service (Business Logic Layer)
	coreService := coreService.NewCoreService("TOBECONTINUED", log, storage)
	coreHandler := core.NewCoreHandler(coreService, log)

	// Создание HTTP обработчика
	httpServer := httpapp.New(log, config.HTTP.Port)
	routes.RegisterRoutes(httpServer.Router, coreHandler)

	return &App{
		HTTPServer: httpServer,
	}
}
