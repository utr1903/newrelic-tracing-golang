package commons

import (
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

func CreateNewRelicAgent() *newrelic.Application {

	log.Info().Msg("Starting New Relic agent...")

	nrapp, err := newrelic.NewApplication(
		newrelic.ConfigEnabled(true),
		newrelic.ConfigAppName("third-go"),
		newrelic.ConfigLicense(os.Getenv("NEWRELIC_LICENSE_KEY")),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigAppLogForwardingMaxSamplesStored(500),
	)

	if err != nil {
		message := "New Relic agent could not be started."
		log.Panic().Msg(message)
		panic(message)
	}

	log.Info().Msg("New Relic agent is started successfully.")

	CreateCustomLogger(nrapp)
	return nrapp
}
