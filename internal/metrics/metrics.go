package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "monitorcontroller_http_response_total",
		Help: "response status for each url",
	}, []string{"url", "status"})

	CrdHttpResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "monitorcontroller_crd_http_response_total",
		Help: "response status totals for each request in a CRD",
	}, []string{"type", "crd", "requestName", "status"})
)
