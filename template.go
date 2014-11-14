package main

import (
	"os"

	"github.com/Shopify/go-base/base.v1"
	"github.com/Shopify/go-dogstatsd/dog"
	log "github.com/Sirupsen/logrus"
)

func main() {
	base.StartPprofServer()
	base.SetupAirbrake()
	base.SetupDatadog("go-base.example.", nil)

	os.Exit(run())
}

func run() int {

	go base.LogPanics(func() {
		if err := dog.Gauge("widgets", 1, nil, 0.001); err != nil {
			log.WithField("error", err).Warn("couldn't open UDP socket -- weird!")
		}
	})

	return 0
}
