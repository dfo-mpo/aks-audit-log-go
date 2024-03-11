package webhook

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jemag/aks-audit-log-go/internal/forwarder"
	"github.com/jemag/aks-audit-log-go/internal/httpclient"
	"github.com/rs/zerolog/log"
)

type WebhookPoster struct {
	forwarderConfiguration *forwarder.ForwarderConfiguration
	httpClient             *httpclient.HttpClientHandler
}

func (w *WebhookPoster) InitConfig(f *forwarder.ForwarderConfiguration) {
	w.forwarderConfiguration = f
	w.httpClient = httpclient.NewHttpClientHandler()
}

func (w *WebhookPoster) PostSyncNoException(url string, contentType string, body string) (*http.Response, error) {
	response, err := w.httpClient.PostAsync(url, contentType, body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (w *WebhookPoster) SendPost(auditEventStr string, mainEventName string, eventNumber int) error {
	retries := 1
	delay := w.forwarderConfiguration.PostRetryIncrementalDelay

	if w.forwarderConfiguration.VerboseLevel > 3 {
		msg := fmt.Sprintf("%s %d > POST", mainEventName, eventNumber)
		log.Info().Msg(msg)
	}

	forwarder.IncreaseSent()

	response, err := w.PostSyncNoException(w.forwarderConfiguration.WebSinkURL, "application/json", auditEventStr)
	if err != nil {
		return err
	}

	status := response.StatusCode == 200 // OK

	for !status && retries <= w.forwarderConfiguration.PostMaxRetries {
		msg := fmt.Sprintf(
			"%s %d > **Error sending POST, retry %d, result: [%d]",
			mainEventName,
			eventNumber,
			retries,
			response.StatusCode,
		)
		log.Error().Msg(msg)

		retries++

		time.Sleep(time.Duration(delay) * time.Millisecond)
		delay += w.forwarderConfiguration.PostRetryIncrementalDelay

		forwarder.IncreaseRetries()

		response, err := w.PostSyncNoException(w.forwarderConfiguration.WebSinkURL, "application/json", auditEventStr)
		if err != nil {
			return err
		}

		status = response.StatusCode == 200 // OK
	}

	if status {
		forwarder.IncreaseSuccesses()
		if w.forwarderConfiguration.VerboseLevel > 3 {
			msg := fmt.Sprintf("%s %d > Post response [%d]", mainEventName, eventNumber, response.StatusCode)
			log.Info().Msg(msg)
		}

		return nil
	} else {
		forwarder.IncreaseErrors()
		msg := fmt.Sprintf("%s %d > **Error post response after max retries, gave up: [%d]", mainEventName, eventNumber, response.StatusCode)
		log.Error().Msg(msg)

		return err
	}
}
