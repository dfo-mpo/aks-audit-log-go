package httpclient

import (
	"net/http"
	"strings"
)

type HttpClientHandler struct {
	client *http.Client
}

func NewHttpClientHandler() *HttpClientHandler {
	return &HttpClientHandler{client: &http.Client{}}
}

func (h *HttpClientHandler) PostAsync(url string, contentType string, body string) (*http.Response, error) {
	return h.client.Post(url, contentType, strings.NewReader(body))
}
