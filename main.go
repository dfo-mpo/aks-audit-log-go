package main

import (
  "fmt"
  "internal/forwarder"
  "internal/webhook"
  "internal/httpclient"
)

func main() {
  fmt.Println("AKS Kuberenetes audit log forwarder from Event Hubs to Agent")
  fmt.Println("Starting Prometheus statistics")

  fmt.Println("Starting Server")
  
  fmt.Println("Testing")

  httpclient.NewHttpClientHandler()

  forwarder.InitServer()

  webhook.InitConfig()

  fmt.Println("Program end")
}
