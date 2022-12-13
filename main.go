package main

import (
  "fmt"
  "github.com/JuanPabloSGU/aks-audit-log-go/internal/forwarder"
)

func main() {
  fmt.Println("AKS Kuberenetes audit log forwarder from Event Hubs to Agent")
  fmt.Println("Starting Prometheus statistics")

  fmt.Println("Starting Server")
  
  fmt.Println("Testing")

  forwarder.InitServer()

  fmt.Println("Program end")
}
