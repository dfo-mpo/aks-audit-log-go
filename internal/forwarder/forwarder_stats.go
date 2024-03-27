package forwarder

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

var (
	sent = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "falco_aks_audit_log_events",
		Help: "Total number of falco events sent",
	})
	errors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "falco_aks_audit_log_events_errors",
		Help: "Total number of falco events sent with error result",
	})
	successes = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "falco_aks_audit_log_events_success",
		Help: "Total number of falco events sent with success result",
	})
	retries = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "falco_aks_audit_log_events_retry",
		Help: "Total number of times falco audit events sent had to be retired",
	})
)

func IncreaseSent() {
	sent.Inc()
}

func IncreaseErrors() {
	errors.Inc()
}

func IncreaseSuccesses() {
	successes.Inc()
}

func IncreaseRetries() {
	retries.Inc()
}

func InitServer() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(sent)
	reg.MustRegister(errors)
	reg.MustRegister(successes)
	reg.MustRegister(retries)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	port := ":9000"
	server := &http.Server{
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal().Msgf("Failed to start server on port %s: %v", port, err)
	}

	log.Debug().Msg("Stopping Server")

	// Wait for server shutdown
	time.Sleep(5 * time.Second)
}
