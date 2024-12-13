package httpclient

import (
	"net/http"
	"strings"
)

type HttpClientHandler struct {
	client *http.Client
}

func NewHttpClientHandler(keepAlive bool) *HttpClientHandler {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DisableKeepAlives = !keepAlive

	return &HttpClientHandler{
		client: &http.Client{Transport: t},
	}
}

func (h *HttpClientHandler) PostAsync(url string, contentType string, body string) (*http.Response, error) {
	return h.client.Post(url, contentType, strings.NewReader(body))
}
