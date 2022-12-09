package main

import (
  "fmt"
  "internal/forwarder"
  "internal/webhook"
)

func main() {
  fmt.Println("AKS Kuberenetes audit log forwarder from Event Hubs to Agent")
  fmt.Println("Starting Prometheus statistics")

  fmt.Println("Starting Server")
  
  fmt.Println("Testing")

  var statistics *forwarder.Statistics = forwarder.Stats

  forwarder.StartServer()

  webhook.InitConfig()

  fmt.Println("Program end")
}
