package common

import (
	"net/http"
	"net/url"
)

type CapturedRequest struct {
	CallCount          uint
	LastRequestHeaders http.Header
	LastQueryParams    url.Values
}

type MockResponse struct {
	UriPath string
	Method  string
	Status  int
	Body    []byte
	Headers http.Header
}

func BuildKey(uriPath, method string) string {
	return uriPath + "::" + method
}

func BuildKeyFromRequest(r *http.Request) string {
	return BuildKey(r.URL.Path, r.Method)
}

type MockResponses map[string]*MockResponse

func (m MockResponses) Get(r *http.Request) (*MockResponse, bool) {
	key := BuildKeyFromRequest(r)
	mock, exists := m[key]
	return mock, exists
}

type CapturedRequests map[string]*CapturedRequest

func (c CapturedRequests) Record(r *http.Request) {
	key := BuildKeyFromRequest(r)
	if _, exists := c[key]; !exists {
		c[key] = &CapturedRequest{}
	}

	c[key].CallCount++
	c[key].LastRequestHeaders = r.Header
	c[key].LastQueryParams = r.URL.Query()
}

func (m CapturedRequests) Get(uriPath, method string) (*CapturedRequest, bool) {
	key := BuildKey(uriPath, method)
	r, exists := m[key]
	return r, exists
}
