package forwarder

import (
  "log"
  "net/http"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

var reg = prometheus.NewRegistry()
var Stats = newStatistics(reg)

type Stat interface {
  IncreaseSent()
  IncreaseErrors()
  IncreaseSuccesses()
  IncreaseRetires()
}

type Statistics struct {
  SentEvents        prometheus.Counter
  Errors            prometheus.Counter
  Successes         prometheus.Counter
  Retries           prometheus.Counter
}

func StartServer() {
  http.Handle("/statistics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func (stat Statistics) IncreaseSent() {
  stat.SentEvents.Add(1)
}
func (stat Statistics) IncreaseErrors() {
  stat.Errors.Add(1)
}
func (stat Statistics) IncreaseSuccesses() {
  stat.Successes.Add(1)
}
func (stat Statistics) IncreaseRetries() {
  stat.Retries.Add(1)
}

func newStatistics(reg prometheus.Registerer) *Statistics  {
  metrics := &Statistics {
    SentEvents: prometheus.NewCounter(prometheus.CounterOpts {
      Name: "falco_aks_audit_log_events",
      Help: "total number of falco events sent",
    }),
    Errors: prometheus.NewCounter(prometheus.CounterOpts {
      Name: "falco_aks_audit_log_events_errors",
      Help: "total number of falco events sent with error result",
    }),
    Successes: prometheus.NewCounter(prometheus.CounterOpts {
      Name: "falco_aks_audit_log_events_success",
      Help: "total number of falco events sent with success result",
    }),
    Retries: prometheus.NewCounter(prometheus.CounterOpts {
      Name: "falco_aks_audit_log_events_retry",
      Help: "total number of times falco audit events sent had to be retired",
    }),
  }

  reg.MustRegister(metrics.SentEvents)
  reg.MustRegister(metrics.Errors)
  reg.MustRegister(metrics.Successes)
  reg.MustRegister(metrics.Retries)

  return metrics
}
