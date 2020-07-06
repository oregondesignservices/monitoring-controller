package v1alpha1

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
)

var (
	HttpClientStatusCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "httpmonitor_client_response_status_total",
		Help: "request statuses for eaach status code",
	}, []string{"url", "status"})

	MonitorRequestByStatus = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "httpmonitor_request_by_status_total",
		Help: "request statuses for each status code",
	}, []string{"type", "crd", "requestName", "status"})
)

func AddHttpClientStatusCounter(url string, status int) {
	HttpClientStatusCounter.WithLabelValues(url, strconv.Itoa(status)).Inc()
}

func AddMonitorRequestByStatus(m *HttpMonitor, req HttpRequest, resp *http.Response) {
	status := 599
	if resp != nil {
		status = resp.StatusCode
	}

	MonitorRequestByStatus.WithLabelValues(
		"HttpMonitor/v1alpha1",
		fmt.Sprintf("%s/%s", m.Namespace, m.Name),
		req.Name,
		strconv.Itoa(status)).Inc()
}
