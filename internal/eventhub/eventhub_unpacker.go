package eventhub

import (
	"encoding/json"
	"fmt"

	"github.com/jemag/aks-audit-log-go/internal/forwarder"
	"github.com/jemag/aks-audit-log-go/internal/webhook"
)

type HubEventUnpacker struct {
	forwarderConfiguration *forwarder.ForwarderConfiguration
	webhookPoster          webhook.WebhookPoster
}

func (h *HubEventUnpacker) InitConfig(f *forwarder.ForwarderConfiguration) {
	h.forwarderConfiguration = f
	h.webhookPoster = webhook.WebhookPoster{}
	h.webhookPoster.InitConfig(f)
}

func (h HubEventUnpacker) Process(eventJObj []byte, mainEventName string) error {
	var event map[string]interface{}
	err := json.Unmarshal(eventJObj, &event)
	if err != nil {
		return err
	}

	for i, record := range event["records"].([]interface{}) {
		record := record.(map[string]interface{})
		auditEventStr := record["properties"].(map[string]interface{})["log"].(string)

		if h.forwarderConfiguration.VerboseLevel > 2 {
			err := h.ConsoleWriteEventSummary(auditEventStr, mainEventName, i)
			if err != nil {
				return err
			}
		}

		err := h.webhookPoster.SendPost(auditEventStr, mainEventName, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h HubEventUnpacker) ConsoleWriteEventSummary(auditEventStr string, mainEventName string, eventNumber int) error {
	var auditEvent map[string]interface{}
	err := json.Unmarshal(([]byte(auditEventStr)), &auditEvent)
	if err != nil {
		return err
	}

	var user, verb, resource, name string

	if userVal, ok := auditEvent["user"].(map[string]interface{})["username"].(string); ok {
		user = userVal
	}

	if verbVal, ok := auditEvent["verb"].(string); ok {
		verb = verbVal
	}

	if objectRef, ok := auditEvent["objectRef"].(map[string]interface{}); ok {
		if resourceVal, ok := objectRef["resource"].(string); ok {
			resource = resourceVal
		}
		if nameVal, ok := objectRef["name"].(string); ok {
			name = nameVal
		}
	}

	fmt.Printf("%s %d > READ audit event: %s %s %s %s",
		mainEventName, eventNumber,
		user,
		verb,
		resource,
		name)

	return nil
}
