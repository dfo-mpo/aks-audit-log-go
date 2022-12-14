package eventhub

import (
  "fmt"
  "encoding/json"
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

func (h HubEventUnpacker) Process(eventJObj []byte, mainEventName string) (bool, error) {
  var event map[string]interface{}
  err := json.Unmarshal(eventJObj, &event)
  if err != nil {
    return false, err
  }

  var results []bool

  ok := true

  for i, record := range event["records"].([]interface{}) {
    record := record.(map[string]interface{})
    auditEventStr := record["properties"].(map[string]interface{})["log"].(string)

    if h.ForwarderConfiguration.VerboseLevel > 2 {
      h.ConsoleWriteEventSummary(auditEventStr, mainEventName, i)
    }

    result, err := h.WebhookPoster.SendPost(auditEventStr, mainEventName, i)
    if err != nil {
      return false, err
    }

    results = append(results, result)
  }

  for _, result := range results {
    ok = ok && result
  }

  return ok, nil
}

func (h HubEventUnpacker) ConsoleWriteEventSummary(auditEventStr string, mainEventName string, eventNumber int) {
  var auditEvent map[string]interface{}
  err := json.Unmarshal([]byte(auditEventStr), &auditEvent)

  if err != nil {
    fmt.Println(err)
  }

  fmt.Printf("%s %d > READ audit event: %s %s %s %s", 
    mainEventName, eventNumber,
    auditEvent["user"].(map[string]interface{})["username"].(string),
    auditEvent["verb"].(string),
    auditEvent["objectRef"].(map[string]interface{})["resource"].(string),
    auditEvent["objectRef"].(map[string]interface{})["name"].(string))
}
