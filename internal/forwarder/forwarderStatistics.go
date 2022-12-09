package forwarder

import (
  "log"
  "net/http"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

type statistics struct {
  sentEvents        prometheus.Counter
  errors            prometheus.Counter
  successes         prometheus.Counter
  retires           prometheus.Counter
}

func InitServer()  {
  reg := prometheus.NewRegistry()
  statistics := newStatistics(reg)

  statistics.sentEvents.Add(1)

  http.Handle("/statistics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

  log.Fatal(http.ListenAndServe(":8080", nil))
}

func newStatistics(reg prometheus.Registerer) *statistics  {
  metrics := &statistics {
    sentEvents: prometheus.NewCounter(prometheus.CounterOpts {
      Name: "falco_aks_audit_log_events",
      Help: "total number of falco events sent",
    }),
    errors: prometheus.NewCounter(prometheus.CounterOpts {
      Name: "falco_aks_audit_log_events_errors",
      Help: "total number of falco events sent with error result",
    }),
    successes: prometheus.NewCounter(prometheus.CounterOpts {
      Name: "falco_aks_audit_log_events_success",
      Help: "total number of falco events sent with success result",
    }),
    retires: prometheus.NewCounter(prometheus.CounterOpts {
      Name: "falco_aks_audit_log_events_retry",
      Help: "total number of times falco audit events sent had to be retired",
    }),
  }

  reg.MustRegister(metrics.sentEvents)
  reg.MustRegister(metrics.errors)
  reg.MustRegister(metrics.successes)
  reg.MustRegister(metrics.retires)

  return metrics
}
