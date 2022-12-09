package forwarder

import (
  "log"
  "net/http"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)


type Statistics struct {
  SentEvents        prometheus.Counter
  Errors            prometheus.Counter
  Successes         prometheus.Counter
  Retries           prometheus.Counter
}

func StartServer() {
  reg := prometheus.NewRegistry()
  statistics := newStatistics(reg)

  statistics.SentEvents.Add(0)

  http.Handle("/statistics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

  log.Fatal(http.ListenAndServe(":8080", nil))
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
