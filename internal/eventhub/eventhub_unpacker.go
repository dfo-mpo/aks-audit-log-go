package eventhub

import (
	"context"
	"encoding/json"

	"github.com/dfo-mpo/aks-audit-log-go/internal/forwarder"
	"github.com/dfo-mpo/aks-audit-log-go/internal/webhook"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
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

type Event struct {
	Records []struct {
		Properties struct {
			Log string `json:"log"`
		} `json:"properties"`
	} `json:"records"`
}

func (h HubEventUnpacker) Process(eventJObj []byte, partitionID string, eventID int64, rateLimiter *rate.Limiter) error {
	err, event := UnmarshallEvent(eventJObj)
	if err != nil {
		return err
	}

	for i, record := range event.Records {
		ctx := context.Background()
		err := rateLimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
		if err != nil {
			return err
		}

		auditEventStr := record.Properties.Log

		log.Debug().Str("partition_id", partitionID).Int64("event_id", eventID).Int("record_id", i).Msgf("%v", record)

		err = h.webhookPoster.SendPost(auditEventStr, partitionID, eventID, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnmarshallEvent(eventJObj []byte) (error, Event) {
	var event Event
	err := json.Unmarshal(eventJObj, &event)
	if err != nil {
		return err, event
	}
	return nil, event
}
