package webhook

import (
  "fmt"
  "net/http"
  "time"
  "github.com/JuanPabloSGU/aks-audit-log-go/internal/forwarder"
  "github.com/JuanPabloSGU/aks-audit-log-go/internal/httpclient"
)

type WebhookPoster struct {
  ForwarderConfiguration    *forwarder.ForwarderConfiguration
  HttpClient                httpclient.IHttpHandler
}

func NewWebhookPoster() *WebhookPoster {
  return &WebhookPoster{}
}

func (w *WebhookPoster) InitConfig() {
  if w.HttpClient == nil {
    w.HttpClient = &httpclient.HttpClientHandler{}
  }
}

func (w *WebhookPoster) PostSyncNoException(url string, content string) (*http.Response, error) {
  response, err := w.HttpClient.PostAsync(url, content)
  if err != nil {
    return nil, err
  }

  return response, nil
}

func (w *WebhookPoster) SendPost(auditEventStr string, mainEventName string, eventNumber int) (bool, error) {
  retries := 1
  delay := w.ForwarderConfiguration.PostRetryIncrementalDelay

  if w.ForwarderConfiguration.VerboseLevel > 3 {
    fmt.Printf("%s %d > POST event to : %s", mainEventName, eventNumber, w.ForwarderConfiguration.WebSinkURL)
  }

  f := forwarder.ForwarderStatistics{}
  f.IncreaseSent()

  response, err := w.PostSyncNoException(w.ForwarderConfiguration.WebSinkURL, auditEventStr)

  status := response.StatusCode == 200 // OK

  for !status && retries <= w.ForwarderConfiguration.PostMaxRetries {
    fmt.Printf("%s %d > **Error sending POST, retry %d, result: [%d] %s", mainEventName, eventNumber, retries, response.StatusCode, response.Body)
    
    retries++

    time.Sleep(time.Duration(delay) * time.Millisecond)
    delay += w.ForwarderConfiguration.PostRetryIncrementalDelay

    f.IncreaseRetries()

    response, err := w.PostSyncNoException(w.ForwarderConfiguration.WebSinkURL, auditEventStr)

    if err != nil {
      fmt.Printf("%s %d > **Error sending POST, retry %d, result: [%d] %s", mainEventName, eventNumber, retries, response.StatusCode, response.Body)
    }

    status = response.StatusCode == 200 // OK
  }

  if status {
    f.IncreaseSuccesses()
    if w.ForwarderConfiguration.VerboseLevel > 3 {
      fmt.Printf("%s %d > Post response OK", mainEventName, eventNumber)
    }

    return true, nil
  } else {
    f.IncreaseErrors()
    fmt.Printf("%s %d > **Error post response after max retries, gave up: [%d] %s", mainEventName, eventNumber, response.StatusCode, response.Body)

    return false, err
  }
}
