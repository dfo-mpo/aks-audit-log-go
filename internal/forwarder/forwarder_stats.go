package forwarder

import (
  "log"
  "net/http"
  "os"
  "os/signal"
  "syscall"
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

type ForwarderStatistics struct {}

func NewForwarderStatistics() *ForwarderStatistics {
  return &ForwarderStatistics{}
}

func (f *ForwarderStatistics) IncreaseSent()  {
  sent.Inc()
}

func (f *ForwarderStatistics) IncreaseErrors() {
  errors.Inc()
}

func (f *ForwarderStatistics) IncreaseSuccesses() {
  successes.Inc()
}

func (f *ForwarderStatistics) IncreaseRetries() {
  retries.Inc()
}

func (f *ForwarderStatistics) StartServer() {
  go InitServer()
}

func InitServer()  {
  reg := prometheus.NewRegistry()
  reg.MustRegister(sent)
  reg.MustRegister(errors)
  reg.MustRegister(successes)
  reg.MustRegister(retries)

  // Create a basic router using serve mux and register the prometheus handler
  mux := http.NewServeMux()
  mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

  // Create a channel to listen for interrupts or signals from OS
  exit := make(chan os.Signal, 1)
  signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
  
  /// Start Server and log errors
  go func() {
    log.Println("Starting Server on PORT : 9000")
    if err := http.ListenAndServe(":9000", mux); err != nil {
      log.Fatal(err)
    }
  }()

  // Wait for an interrupt or signal from OS
  <-exit
  log.Println("Stopping Server")

  // Wait for server shutdown
  time.Sleep(5 * time.Second)
}
