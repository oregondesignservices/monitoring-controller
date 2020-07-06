package v1alpha1

var KnownRunners map[string]*HttpMonitorRunner

func init() {
	KnownRunners = make(map[string]*HttpMonitorRunner)
}
