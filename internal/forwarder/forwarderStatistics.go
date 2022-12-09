package forwarder

import (
  "github.com/prometheus/client_golang/prometheus"
)

type statistics struct {
  sentEvents        prometheus.Counter
  errors            prometheus.Counter
  successes         prometheus.Counter
  retires           prometheus.Counter
}

func NewMetrics(reg prometheus.Registerer) *statistics  {
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
