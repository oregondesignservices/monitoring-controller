package httpmonitor

import (
	monitoringraisingthefloororgv1alpha1 "github.com/oregondesignservices/monitoring-controller/api/v1alpha1"
)

var KnownHttpMonitors map[string]*monitoringraisingthefloororgv1alpha1.HttpMonitor

func init() {
	KnownHttpMonitors = make(map[string]*monitoringraisingthefloororgv1alpha1.HttpMonitor)
}
