/*
Copyright 2020 Raising The Floor.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	monitoringraisingthefloororgv1alpha1 "github.com/oregondesignservices/monitoring-controller/api/v1alpha1"
	"github.com/oregondesignservices/monitoring-controller/controllers"
	"github.com/oregondesignservices/monitoring-controller/internal/conf"
	"github.com/urfave/cli/v2"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = monitoringraisingthefloororgv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	app := &cli.App{
		Name:   "monitoring-controller",
		Usage:  "periodically make requests to monitor service health",
		Action: run,
		Flags:  conf.Flags,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func run(c *cli.Context) error {
	err := conf.GlobalConfig.UpdateFromCli(c)
	if err != nil {
		return err
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      conf.GlobalConfig.MetricsAddr,
		Port:                    9443,
		LeaderElection:          conf.GlobalConfig.EnableLeaderElection,
		LeaderElectionID:        "monitoring-controller.raisingthefloor.org",
		LeaderElectionNamespace: conf.GlobalConfig.Namespace,
		//HealthProbeBindAddress:  ":8081",
		//ReadinessEndpointName:   "/ready",
		//LivenessEndpointName:    "/alive",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.HttpMonitorReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("HttpMonitor"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HttpMonitor")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
	return nil
}
