package eventhub

import (
  "github.com/JuanPabloSGU/aks-audit-log-go/internal/forwarder"
  "github.com/JuanPabloSGU/aks-audit-log-go/internal/webhook"
)

type HubEventUnpacker struct {
  ForwarderConfiguration  *forwarder.ForwarderConfiguration
  WebhookPoster           webhook.WebhookPoster 
}

func (h HubEventUnpacker) InitConfig(f *forwarder.ForwarderConfiguration) {
  h.WebhookPoster = webhook.WebhookPoster{ForwarderConfiguration: f,}
  h.WebhookPoster.InitConfig()
}

func (h HubEventUnpacker) Process(eventMap []byte, mainEventName string) (bool, error) {
  return true, nil
}

