package app

import (
	"log/slog"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	httpapp "github.com/shameoff/more-than-trip/core/internal/app/http"
	"github.com/shameoff/more-than-trip/core/internal/config"
	"github.com/shameoff/more-than-trip/core/internal/http-server/core"
	routes "github.com/shameoff/more-than-trip/core/internal/http-server/core"
	"github.com/shameoff/more-than-trip/core/internal/lib/converter"
	"github.com/shameoff/more-than-trip/core/internal/lib/logger/sl"
	coreService "github.com/shameoff/more-than-trip/core/internal/services/core"
	"github.com/shameoff/more-than-trip/core/internal/storage/postgres"
	s3Storage "github.com/shameoff/more-than-trip/core/internal/storage/s3"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs" // docs is generated by Swag CLI, you have to import it.

	httpSwagger "github.com/swaggo/http-swagger/v2"
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

	// Init S3 connection
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"), // Укажите регион
	})
	if err != nil {
		log.Error("failed to create S3 session", sl.Err(err))
		panic(err)
	}
	s3Client := s3.New(sess)
	s3PhotoService := s3Storage.NewS3Storage(config.S3.PhotosBucket, s3Client)

	// Init core service (Business Logic Layer)
	coreService := coreService.NewCoreService("TOBECONTINUED", log, storage, s3PhotoService)
	coreHandler := core.NewCoreHandler(coreService, log)

	// Создание HTTP обработчика
	httpServer := httpapp.New(log, config.HTTP.Port)
	routes.RegisterRoutes(httpServer.Router, coreHandler)

	// Init swagger
	httpServer.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/swagger/doc.json"), //The url pointing to API definition
	))

	http.ListenAndServe(":1323", httpServer.Router)
	return &App{
		HTTPServer: httpServer,
	}
}
