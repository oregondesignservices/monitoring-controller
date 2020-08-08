// a simple server to capture HTTP monitor requests
package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

var logger *zap.SugaredLogger

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = l.Sugar()
}

type RequestDatum struct {
	CallCount          uint
	LastRequestHeaders http.Header
	LastQueryParams    url.Values
}

type ConfiguredResponse struct {
	Status  int
	Body    []byte
	Headers http.Header
}

func buildKey(r *http.Request) string {
	return r.URL.Path + "::" + r.Method
}

type ConfiguredResponses map[string]ConfiguredResponse

type RequestData map[string]*RequestDatum

var howToRespond = ConfiguredResponses{}
var capturedRequests = RequestData{}

func main() {
	logger.Info("starting...")
	defer func() {
		logger.Info("shutting down")
	}()

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/reset", ResetHandler)
	http.HandleFunc("/data", DataHandler)
	http.HandleFunc("/health", HealthHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("ListenAndServe returned an error", zap.Error(err))
	}
}

// Returns responses
func RootHandler(w http.ResponseWriter, r *http.Request) {
	key := buildKey(r)
	entry := logger.With(zap.String("key", key))

	// Record the request
	if _, exists := capturedRequests[key]; !exists {
		capturedRequests[key] = &RequestDatum{}
	}

	capturedRequests[key].CallCount++
	capturedRequests[key].LastRequestHeaders = r.Header
	capturedRequests[key].LastQueryParams = r.URL.Query()

	// Figure out how to respond.
	response, exists := howToRespond[key]
	if !exists {
		entry.Info("I have not been told how to respond")
		w.WriteHeader(http.StatusNotFound)
		_, err := fmt.Fprintf(w, "I have not been told how to respond to '%s'", key)
		if err != nil {
			entry.Error("failed to encode response", zap.Error(err))
		}
		return
	}

	entry.Info("I know how how to respond")
	for key, val := range response.Headers {
		for _, v := range val {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(response.Status)
	_, err := w.Write(response.Body)
	if err != nil {
		entry.Error("failed to encode response", zap.Error(err))
	}
}

// Reset all data
func ResetHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("resetting request and response data")
	howToRespond = ConfiguredResponses{}
	capturedRequests = RequestData{}
}

// Get the data
func DataHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("dumping captured requests")
	err := jsoniter.NewEncoder(w).Encode(capturedRequests)
	if err != nil {
		logger.Error("failed to encode response", zap.Error(err))
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// does nothing. just need a 200 handler
}
