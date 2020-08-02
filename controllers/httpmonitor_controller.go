/*
Copyright 2020 Raising the Floor - International

Licensed under the New BSD license. You may not use this file except in
compliance with this License.

You may obtain a copy of the License at
https://github.com/GPII/universal/blob/master/LICENSE.txt

The R&D leading to these results received funding from the:
* Rehabilitation Services Administration, US Dept. of Education under
  grant H421A150006 (APCP)
* National Institute on Disability, Independent Living, and
  Rehabilitation Research (NIDILRR)
* Administration for Independent Living & Dept. of Education under grants
  H133E080022 (RERC-IT) and H133E130028/90RE5003-01-00 (UIITA-RERC)
* European Union's Seventh Framework Programme (FP7/2007-2013) grant
  agreement nos. 289016 (Cloud4all) and 610510 (Prosperity4All)
* William and Flora Hewlett Foundation
* Ontario Ministry of Research and Innovation
* Canadian Foundation for Innovation
* Adobe Foundation
* Consumer Electronics Association Foundation
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/oregondesignservices/monitoring-controller/internal/conf"
	"github.com/oregondesignservices/monitoring-controller/internal/metrics"
	runnverv1alpha1 "github.com/oregondesignservices/monitoring-controller/internal/runner/v1alpha1"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"

	monitoringraisingthefloororgv1alpha1 "github.com/oregondesignservices/monitoring-controller/api/v1alpha1"
)

// HttpMonitorReconciler reconciles a HttpMonitor object
type HttpMonitorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=monitoring.raisingthefloor.org,resources=httpmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.raisingthefloor.org,resources=httpmonitors/status,verbs=get;update;patch

func (r *HttpMonitorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	instance := &monitoringraisingthefloororgv1alpha1.HttpMonitor{}
	ctx := context.Background()
	logger := r.Log.WithValues("httpmonitor", req.NamespacedName, "key", req.NamespacedName.String())

	runnerKey := req.NamespacedName.String()
	knownRunner, runnerExists := runnverv1alpha1.KnownRunners[runnerKey]

	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			removeKnownHttpCrdGauge(logger, req.Namespace, req.Name)
			// Object not found. See if we need to stop a monitor
			if runnerExists {
				logger.Info("removing monitor")
				knownRunner.Stop()
				delete(runnverv1alpha1.KnownRunners, runnerKey)
			}
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	for k, v := range conf.GlobalConfig.GlobalRequestVars {
		instance.Spec.Globals[k] = v
	}

	if !runnerExists {
		logger.Info("detected a new http monitor")
	} else {
		// If the resource version is the same, we have nothing to do. We know about the exact object.
		if instance.GetResourceVersion() == knownRunner.GetResourceVersion() {
			logger.V(3).Info("received a known http monitor with no changes")
			return reconcile.Result{}, nil
		} else {
			logger.Info("detected http monitor changes")
			knownRunner.Stop()
			removeKnownHttpCrdGauge(logger, req.Namespace, req.Name)
		}
	}

	recordKnownHttpCrdGauge(instance)

	// At this point, we need to store the http monitor and restart its worker routine
	newRunner := runnverv1alpha1.NewHttpMonitorRunner(instance)
	runnverv1alpha1.KnownRunners[runnerKey] = newRunner
	newRunner.Start()

	return ctrl.Result{}, nil
}

func labelPairsToLabels(pairs []*dto.LabelPair) prometheus.Labels {
	m := prometheus.Labels{}

	for _, pair := range pairs {
		m[*pair.Name] = *pair.Value
	}

	return m
}

// We need to remove the existing gauge so we can update its details
func removeKnownHttpCrdGauge(logger logr.Logger, namespace, name string) {
	ch := make(chan prometheus.Metric)
	go func() {
		metrics.KnownHttpCrdGauge.Collect(ch)
		close(ch)
	}()

	var toDelete []prometheus.Labels
	for m := range ch {
		pb := &dto.Metric{}
		err := m.Write(pb)
		if err != nil {
			ctrl.Log.Error(err, "failed to decode metric")
		}
		labels := labelPairsToLabels(pb.GetLabel())
		if labels["namespace"] == namespace && labels["name"] == name {
			toDelete = append(toDelete, labels)
		}
	}

	for _, labelsToDelete := range toDelete {
		logger.Info("deleting existing metric in KnownHttpCrdGauge", "labels", labelsToDelete)
		metrics.KnownHttpCrdGauge.Delete(labelsToDelete)
	}
}

func recordKnownHttpCrdGauge(crd *monitoringraisingthefloororgv1alpha1.HttpMonitor) {
	metrics.KnownHttpCrdGauge.With(prometheus.Labels{
		"namespace":            crd.Namespace,
		"name":                 crd.Name,
		"num_requests":         strconv.Itoa(len(crd.Spec.Requests)),
		"num_cleanup_requests": strconv.Itoa(len(crd.Spec.Cleanup)),
		"period":               crd.Spec.Period.Duration.String(),
		"num_globals":          strconv.Itoa(len(crd.Spec.Globals)),
	}).Set(1)
}

func (r *HttpMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringraisingthefloororgv1alpha1.HttpMonitor{}).
		Complete(r)
}
