package main

import (
	"os"
	"strings"

	"github.com/jemag/aks-audit-log-go/internal/eventhub"
	"github.com/jemag/aks-audit-log-go/internal/forwarder"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	logLevelStr := strings.ToLower(os.Getenv("LOGLEVEL"))
	var logLevel zerolog.Level = zerolog.InfoLevel
	switch logLevelStr {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	case "fatal":
		logLevel = zerolog.FatalLevel
	case "panic":
		logLevel = zerolog.PanicLevel
	case "trace":
		logLevel = zerolog.TraceLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	log.Debug().Msg("AKS Kuberenetes audit log forwarder from Event Hubs to Agent")

	log.Debug().Msg("Starting Prometheus statistics")

	log.Debug().Msg("Starting Server")

	go forwarder.InitServer()

	log.Debug().Msg("Testing")

	eventhub.Run()

	log.Debug().Msg("Program end")
}
