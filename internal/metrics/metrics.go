package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	HttpResponseCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "monitor_http_response_total",
		Help: "response status for each url",
	}, []string{"url", "status"})

	CrdHttpResponseCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "monitor_crd_http_response_total",
		Help: "response status totals for each request in a CRD",
	}, []string{"type", "crd", "requestName", "status"})

	KnownHttpCrdGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "monitor_http_crd_details",
		Help: "details for HttpMonitor CRDs",
	}, []string{"namespace", "name", "num_requests", "num_cleanup_requests", "period", "num_globals"})
)

func init() {
	metrics.Registry.MustRegister(
		HttpResponseCounter,
		CrdHttpResponseCounter,
		KnownHttpCrdGauge)
}
