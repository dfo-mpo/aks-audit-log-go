package httpclient

import (
  "net/http"
)

type HttpClientHandler struct {
  client *http.Client
}

func NewHttpClientHandler() *HttpClientHandler {
  return &HttpClientHandler{client: &http.Client{}}
}

func (h *HttpClientHandler) GetAsync(url string) (*http.Response, error) {
  return h.client.Get(url)
}

func (h *HttpClientHandler) PostAsync(url string, content string) (*http.Response, error)  {
 return h.client.Post(url, content, nil)
}
