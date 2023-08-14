package forwarder

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	log.Fatal(http.ListenAndServe(":9000", nil))

	log.Println("Stopping Server")

	// Wait for server shutdown
	time.Sleep(5 * time.Second)
}
