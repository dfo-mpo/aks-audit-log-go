package main

import (
	"flag"

	"github.com/jemag/aks-audit-log-go/internal/eventhub"
	"github.com/jemag/aks-audit-log-go/internal/forwarder"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Default level for this example is info, unless debug flag is present
	debug := flag.Bool("debug", false, "sets log level to debug")
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Debug().Msg("AKS Kuberenetes audit log forwarder from Event Hubs to Agent")

	log.Debug().Msg("Starting Prometheus statistics")

	log.Debug().Msg("Starting Server")

	go forwarder.InitServer()

	log.Debug().Msg("Testing")

	eventhub.Run()

	log.Debug().Msg("Program end")
}
