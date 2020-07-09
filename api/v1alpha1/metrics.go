package v1alpha1

import (
	"fmt"
	"github.com/oregondesignservices/monitoring-controller/internal/metrics"
	"net/http"
	"strconv"
)

func HandleMetrics(m *HttpMonitor, req HttpRequest, resp *http.Response) {
	status := 599
	if resp != nil {
		status = resp.StatusCode
	}
	stringStatus := strconv.Itoa(status)

	metrics.HttpResponseCounter.WithLabelValues(req.Url, stringStatus).Inc()

	metrics.CrdHttpResponseCounter.WithLabelValues(
		"HttpMonitor/v1alpha1",
		fmt.Sprintf("%s/%s", m.Namespace, m.Name),
		req.Name,
		stringStatus).Inc()
}
