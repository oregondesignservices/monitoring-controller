package httpclient

import (
	"net/http"
	"time"
)

var httpClient *http.Client

func Initialize(timeout time.Duration) {
	httpClient = &http.Client{
		Timeout: timeout,
	}
}

func GetClient() *http.Client {
	return httpClient
}
