package v1alpha1

import (
	monitoringraisingthefloororgv1alpha1 "github.com/oregondesignservices/monitoring-controller/api/v1alpha1"
	"sync"
	"time"
)

type HttpMonitorRunner struct {
	*monitoringraisingthefloororgv1alpha1.HttpMonitor
	ticker  *time.Ticker
	stopped *sync.WaitGroup
}

func NewHttpMonitorRunner(m *monitoringraisingthefloororgv1alpha1.HttpMonitor) *HttpMonitorRunner {
	return &HttpMonitorRunner{HttpMonitor: m}
}

func (h *HttpMonitorRunner) Start() {
	if h.ticker != nil {
		panic("tried to start an already started HttpMonitor")
	}

	h.ticker = time.NewTicker(h.Spec.Period.Duration)
	h.stopped = &sync.WaitGroup{}
	h.stopped.Add(1)
	go func() {
		defer h.stopped.Done()
		for _ = range h.ticker.C {
			h.Execute()
		}
	}()
}

func (h *HttpMonitorRunner) Stop() {
	h.ticker.Stop()
	h.stopped.Wait()
}
