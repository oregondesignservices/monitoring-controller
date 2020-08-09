// a simple server to capture HTTP monitor requests
package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/oregondesignservices/monitoring-controller/tester/common"
	"go.uber.org/zap"
	"net/http"
)

var logger *zap.SugaredLogger

func init() {
	logger = common.Logger()
}

var howToRespond = common.MockResponses{}
var capturedRequests = common.CapturedRequests{}

func main() {
	logger.Info("starting...")
	defer func() {
		logger.Info("shutting down")
	}()

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/reset", ResetHandler)
	http.HandleFunc("/set", SetHandler)
	http.HandleFunc("/data", DataHandler)
	http.HandleFunc("/health", HealthHandler)
	if err := http.ListenAndServe(":80", nil); err != nil {
		logger.Error("ListenAndServe returned an error", zap.Error(err))
	}
}

// Returns responses
func RootHandler(w http.ResponseWriter, r *http.Request) {
	entry := logger.With(zap.String("uriPath", r.URL.Path), zap.String("method", r.Method))

	// Record the request
	capturedRequests.Record(r)

	// Figure out how to respond.
	response, exists := howToRespond.Get(r)
	if !exists {
		entry.Info("I have not been told how to respond")
		w.WriteHeader(http.StatusNotFound)
		_, err := fmt.Fprintf(w, "I have not been told how to respond to '%s'", common.BuildKeyFromRequest(r))
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

// configure how to respond
func SetHandler(w http.ResponseWriter, r *http.Request) {
	config := &common.MockResponse{}
	err := jsoniter.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		logger.Error("failed to decode response", zap.Error(err))
		w.WriteHeader(500)
		return
	}
	key := common.BuildKey(config.UriPath, config.Method)
	logger.Info("I have been told how to respond", zap.String("key", key))
	howToRespond[key] = config
}

// Reset all data
func ResetHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("resetting request and response data")
	howToRespond = common.MockResponses{}
	capturedRequests = common.CapturedRequests{}
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
	w.WriteHeader(http.StatusOK)
}
