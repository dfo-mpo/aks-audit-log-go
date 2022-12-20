package eventhub

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/JuanPabloSGU/aks-audit-log-go/internal/forwarder"
	"github.com/JuanPabloSGU/aks-audit-log-go/internal/webhook"
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

func (h HubEventUnpacker) Process(eventJObj []byte, mainEventName string) (bool, error) {
	var event map[string]interface{}
	err := json.Unmarshal(eventJObj, &event)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	var results []bool

	ok := true

	for i, record := range event["records"].([]interface{}) {
		record := record.(map[string]interface{})
		auditEventStr := record["properties"].(map[string]interface{})["log"].(string)

		if h.forwarderConfiguration.VerboseLevel > 2 {
			// h.ConsoleWriteEventSummary(auditEventStr, mainEventName, i)
		}

		result, err := h.webhookPoster.SendPost(auditEventStr, mainEventName, i)
		if err != nil {
			log.Fatalln(err)
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
	err := json.Unmarshal(([]byte(auditEventStr)), &auditEvent)

	if err != nil {
		log.Fatalln(err)
		fmt.Println(err)
	}

	fmt.Printf("%s %d > READ audit event: %s %s %s %s",
		mainEventName, eventNumber,
		auditEvent["user"].(map[string]interface{})["username"].(string),
		auditEvent["verb"].(string),
		auditEvent["objectRef"].(map[string]interface{})["resource"].(string),
		auditEvent["objectRef"].(map[string]interface{})["name"].(string))
}
