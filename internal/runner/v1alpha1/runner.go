package v1alpha1

import (
	monitoringraisingthefloororgv1alpha1 "github.com/oregondesignservices/monitoring-controller/api/v1alpha1"
	"time"
)

type HttpMonitorRunner struct {
	*monitoringraisingthefloororgv1alpha1.HttpMonitor
	ticker *time.Ticker
	closer chan bool
}

func NewHttpMonitorRunner(m *monitoringraisingthefloororgv1alpha1.HttpMonitor) *HttpMonitorRunner {
	return &HttpMonitorRunner{HttpMonitor: m}
}

func (h *HttpMonitorRunner) Start() {
	if h.ticker != nil {
		panic("tried to start an already started HttpMonitor")
	}

	h.ticker = time.NewTicker(h.Spec.Period.Duration)
	h.closer = make(chan bool)
	go func() {
		for {
			select {
			case <-h.ticker.C:
				h.Execute()
			case <-h.closer:
				return
			}
		}
	}()
}

func (h *HttpMonitorRunner) Stop() {
	// Stop does not close the channel, so the closer channel handles that.
	h.closer <- true
	h.ticker.Stop()
}
