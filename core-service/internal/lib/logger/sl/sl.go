package sl

import (
	"log/slog"
	"os"
)

const (
	envLocal   = "local"
	envDev     = "dev"
	envProd    = "prod"
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

func New(env string) *slog.Logger {
	var log *slog.Logger

	opts := slog.HandlerOptions{
		// Level reports the minimum record level that will be logged
		Level: LevelTrace,
		// ReplaceAttr in current case gives new levels their names in logs
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	}

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &opts),
		)
	case envDev:
		opts.Level = slog.LevelDebug
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &opts),
		)
	case envProd:
		opts.Level = slog.LevelInfo
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &opts),
		)
	}
	return log
}
