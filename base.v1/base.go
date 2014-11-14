package base

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	dog "github.com/Shopify/go-dogstatsd/dog"
	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/airbrake"
	"github.com/tobi/airbrake-go"
)

const (
	PprofEndpoint       = "PPROF_ENDPOINT"
	AirbrakeEndpoint    = "AIRBRAKE_ENDPOINT"
	AirbrakeAPIKey      = "AIRBRAKE_API_KEY"
	AirbrakeEnvironment = "AIRBRAKE_ENVIRONMENT"
	StatsdEndpoint      = "STATSD_ENDPOINT"
)

func LogPanics(f func()) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("Recovered from panic(%+v)", r)
			log.WithField("error", err).Panicf("terminating due to panic: %s", err.Error())
		}
	}()
	f()
}

func SetupAirbrake() {
	ep := os.Getenv(AirbrakeEndpoint)
	key := os.Getenv(AirbrakeAPIKey)
	env := os.Getenv(AirbrakeEnvironment)

	var missing []string
	if ep == "" {
		missing = append(missing, AirbrakeEndpoint)
	}
	if key == "" {
		missing = append(missing, AirbrakeAPIKey)
	}
	if key == "" {
		missing = append(missing, AirbrakeEnvironment)
	}

	if len(missing) > 0 {
		log.WithField("envvar", missing).Warn("airbrake configuration key(s) not provided; can't configure error reporting")
		return
	}

	airbrake.Endpoint = ep
	airbrake.ApiKey = key
	airbrake.Environment = env

	log.AddHook(new(logrus_airbrake.AirbrakeHook))
	log.WithFields(log.Fields{"endpoint": ep, "environment": env}).Info("successfully configured airbrake")
}

func StartPprofServer() {
	ep := os.Getenv(PprofEndpoint)
	if ep == "" {
		ep = "localhost:6060"
		log.WithFields(log.Fields{"envvar": PprofEndpoint, "default": ep}).Warn("envvar not provided; using default value")
	}

	go LogPanics(func() {
		log.WithField("bind-addr", ep).Info("starting pprof server")
		err := http.ListenAndServe(ep, nil)
		log.WithField("error", err).Error("pprof HTTP server shut down")
	})
}

func SetupDatadog(namespace string, tags []string) {
	ep := os.Getenv(StatsdEndpoint)
	if ep == "" {
		log.WithField("envvar", StatsdEndpoint).Warn("envvar missing; won't submit to statsd")
		return
	}

	if err := dog.Configure(ep, namespace, tags); err != nil {
		log.WithField("error", err).Error("failed to configure statsd/datadog")
	}
}
