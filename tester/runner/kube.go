package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func newKubeClientset() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, err
}

func KubectlApply(path string) error {
	return RunInteractiveCommandNoPipes("kubectl",
		"--kubeconfig", kubeconfig,
		"apply", "-f", path)
}

func KubectlDelete(path string) error {
	return RunInteractiveCommandNoPipes("kubectl",
		"--kubeconfig", kubeconfig,
		"delete", "-f", path)
}
