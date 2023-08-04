package main

import (
	"fmt"
	"github.com/jemag/aks-audit-log-go/internal/eventhub"
	"github.com/jemag/aks-audit-log-go/internal/forwarder"
)

func main() {
	fmt.Println("AKS Kuberenetes audit log forwarder from Event Hubs to Agent")
	fmt.Println("Starting Prometheus statistics")

	fmt.Println("Starting Server")

	go forwarder.InitServer()

	fmt.Println("Testing")

	eventhub.Run()

	fmt.Println("Program end")
}
