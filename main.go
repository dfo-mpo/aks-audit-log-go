package main

import (
  "fmt"
  //"internal/eventhub"
  "internal/forwarder"
  //"internal/webhook"

  "github.com/prometheus/client_golang/prometheus"
)

func main() {
  fmt.Println("AKS Kuberenetes audit log forwarder from Event Hubs to Agent")
  fmt.Println("Starting Prometheus statistics")

  fmt.Println("Starting Server")
  
  fmt.Println("Testing")

  //eventhub.Foo()
  //eventhub.Bar()

  //webhook.WebP()
  //webhook.Web()

  //forwarder.Bar()
  //forwarder.Foo()

  fmt.Printf(forwarder.GetEventHubConnection())

  reg := prometheus.NewRegistry()
  metrics := NewMetrics(reg)

  fmt.Println("Program end")
}
