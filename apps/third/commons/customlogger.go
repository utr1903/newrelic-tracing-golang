package commons

import (
	"context"
	"os"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrzerolog"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CustomLogger struct {
	Nrapp  *newrelic.Application
	Logger zerolog.Logger
}

var logger *CustomLogger

func CreateCustomLogger(
	nrapp *newrelic.Application,
) {

	log.Info().Msg("Initializing custom logger...")

	logger = &CustomLogger{
		Nrapp:  nrapp,
		Logger: zerolog.New(os.Stdout),
	}

	logger.Logger.Info().Msg("Custom logger is initialized successfully.")
}

func Log(
	logLevel zerolog.Level,
	message string,
) {

	nrLogger := logger.Logger.Hook(nrzerolog.NewRelicHook{
		App: logger.Nrapp,
	})

	if logger == nil {
		panic("Custom logger is not initiated.")
	} else {

		switch logLevel {
		case zerolog.ErrorLevel:
			nrLogger.Error().Msg(message)
		case zerolog.PanicLevel:
			nrLogger.Panic().Msg(message)
		default:
			nrLogger.Info().Msg(message)
		}
	}
}

func LogWithContext(
	txn *newrelic.Transaction,
	logLevel zerolog.Level,
	message string,
) {

	ctx := newrelic.NewContext(context.Background(), txn)

	nrLogger := logger.Logger.Hook(nrzerolog.NewRelicHook{
		App:     logger.Nrapp,
		Context: ctx,
	})

	if logger == nil {
		panic("Custom logger is not initiated.")
	} else {

		switch logLevel {
		case zerolog.ErrorLevel:
			nrLogger.Error().Msg(message)
		case zerolog.PanicLevel:
			nrLogger.Panic().Msg(message)
		default:
			nrLogger.Info().Msg(message)
		}
	}
}
