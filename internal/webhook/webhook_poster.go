package webhook

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JuanPabloSGU/aks-audit-log-go/internal/forwarder"
	"github.com/JuanPabloSGU/aks-audit-log-go/internal/httpclient"
)

type WebhookPoster struct {
  forwarderConfiguration    *forwarder.ForwarderConfiguration
  httpClient                *httpclient.HttpClientHandler
}

func NewWebhookPoster() *WebhookPoster {
  return &WebhookPoster{}
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

func (w *WebhookPoster) SendPost(auditEventStr string, mainEventName string, eventNumber int) (bool, error) {
  retries := 1
  delay := w.forwarderConfiguration.PostRetryIncrementalDelay

  if w.forwarderConfiguration.VerboseLevel > 3 {
    fmt.Printf("%s %d > POST event to : %s", mainEventName, eventNumber, w.forwarderConfiguration.WebSinkURL)
  }

  f := forwarder.ForwarderStatistics{}
  f.IncreaseSent()

  response, err := w.PostSyncNoException(w.forwarderConfiguration.WebSinkURL, "application/json", auditEventStr)

  if err != nil {
    log.Fatalln(err)
  }

  status := response.StatusCode == 200 // OK


  for !status && retries <= w.forwarderConfiguration.PostMaxRetries {
    fmt.Printf("%s %d > **Error sending POST, retry %d, result: [%d] %s", mainEventName, eventNumber, retries, response.StatusCode, response.Body)
    
    retries++

    time.Sleep(time.Duration(delay) * time.Millisecond)
    delay += w.forwarderConfiguration.PostRetryIncrementalDelay

    f.IncreaseRetries()

    response, err := w.PostSyncNoException(w.forwarderConfiguration.WebSinkURL, "application/json", auditEventStr)

    if err != nil {
      fmt.Printf("%s %d > **Error sending POST, retry %d, result: [%d] %s", mainEventName, eventNumber, retries, response.StatusCode, response.Body)
    }

    status = response.StatusCode == 200 // OK
  }

  if status {
    f.IncreaseSuccesses()
    if w.forwarderConfiguration.VerboseLevel > 3 {
      fmt.Printf("%s %d > Post response OK", mainEventName, eventNumber)
    }

    return true, nil
  } else {
    f.IncreaseErrors()
    fmt.Printf("%s %d > **Error post response after max retries, gave up: [%d] %s", mainEventName, eventNumber, response.StatusCode, response.Body)

    return false, err
  }
}
