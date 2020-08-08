package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func GetSecret(namespace, name string) (*corev1.Secret, error) {
	return kubeClientset.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
}

func KubectlApply(namespace, path string) error {
	return RunInteractiveCommandNoPipes("kubectl", "apply", "-n", namespace, "-f", path)
}

func KubectlDelete(namespace, path string) error {
	return RunInteractiveCommandNoPipes("kubectl", "delete", "-n", namespace, "-f", path)
}
