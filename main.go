package main

import (
  "fmt"
  "github.com/JuanPabloSGU/aks-audit-log-go/internal/forwarder"
  "github.com/JuanPabloSGU/aks-audit-log-go/internal/webhook"
  "github.com/JuanPabloSGU/aks-audit-log-go/internal/httpclient"
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
