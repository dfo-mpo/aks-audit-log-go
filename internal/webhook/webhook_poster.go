package webhook

import (
	"net/http"
	"time"

	"github.com/dfo-mpo/aks-audit-log-go/internal/forwarder"
	"github.com/dfo-mpo/aks-audit-log-go/internal/httpclient"
	"github.com/rs/zerolog/log"
)

type WebhookPoster struct {
	forwarderConfiguration *forwarder.ForwarderConfiguration
	httpClient             *httpclient.HttpClientHandler
}

func (w *WebhookPoster) InitConfig(f *forwarder.ForwarderConfiguration) {
	w.forwarderConfiguration = f
	w.httpClient = httpclient.NewHttpClientHandler(w.forwarderConfiguration.KeepAlive)
}

func (w *WebhookPoster) PostSyncNoException(url string, contentType string, body string, keepAlive bool) (*http.Response, error) {
	response, err := w.httpClient.PostAsync(url, contentType, body)
	if err != nil {
		return nil, err
	}

	if !keepAlive {
		response.Body.Close()
	}

	return response, nil
}

func (w *WebhookPoster) SendPost(auditEventStr string, partitionID string, eventID int64, recordID int) error {
	retries := 0
	delay := w.forwarderConfiguration.PostRetryIncrementalDelay

	log.Debug().Str("partition_id", partitionID).Int64("event_id", eventID).Int("record_id", recordID).Msg("POST")

	forwarder.IncreaseSent()

	response, err := w.PostSyncNoException(w.forwarderConfiguration.WebSinkURL, "application/json", auditEventStr, w.forwarderConfiguration.KeepAlive)
	if err != nil {
		return err
	}

	status := response.StatusCode == 200 // OK

	for !status && retries <= w.forwarderConfiguration.PostMaxRetries {
		log.Error().Str("partition_id", partitionID).Int64("event_id", eventID).Int("record_id", recordID).Int("retries", retries).Int("status_code", response.StatusCode).Msg("POST unseccessful")

		retries++

		time.Sleep(time.Duration(delay) * time.Millisecond)
		delay += w.forwarderConfiguration.PostRetryIncrementalDelay

		forwarder.IncreaseRetries()

		response, err := w.PostSyncNoException(w.forwarderConfiguration.WebSinkURL, "application/json", auditEventStr, w.forwarderConfiguration.KeepAlive)
		if err != nil {
			return err
		}

		status = response.StatusCode == 200 // OK
	}

	if status {
		forwarder.IncreaseSuccesses()
		log.Debug().Str("partition_id", partitionID).Int64("event_id", eventID).Int("record_id", recordID).Int("retries", retries).Int("status_code", response.StatusCode).Msg("POST successful")
		return nil
	} else {
		forwarder.IncreaseErrors()
		log.Error().Str("partition_id", partitionID).Int64("event_id", eventID).Int("record_id", recordID).Int("retries", retries).Int("status_code", response.StatusCode).Msg("POST unsuccessful, max retries")
		return err
	}
}
