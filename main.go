package main

import (
  "fmt"
  //"internal/eventhub"
  "internal/forwarder"
  //"internal/webhook"
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

  forwarder.InitConfig()

  fmt.Println("Program end")
}
