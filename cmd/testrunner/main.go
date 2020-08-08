package main

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	httpClient = http.Client{
		Timeout: 29 * time.Second,
	}
	kubeClientset *kubernetes.Clientset
)

const (
	metricsUrl  = "http://localhost:9090/metrics"
	kubeconfig  = "kind-kubeconfig.yaml"
	testdataDir = "cmd/testrunner/testdata"
)

func main() {
	var err error
	kubeClientset, err = newKubeClientset()
	if err != nil {
		panic("failed to construct kubernetes clientset: " + err.Error())
	}

	runners := []TestRunner{
		// Wait for controller to wake up so we start real tests
		&WaitForController{},
	}

	successCount := 0
	failureCount := 0

	for _, runner := range runners {
		err := runOne(runner)
		if err == nil {
			successCount++
		} else {
			failureCount++
		}
	}

	fmt.Println("-------------------------------------------------")
	fmt.Printf("successful tests: %d\n", successCount)
	fmt.Printf("failed tests: %d\n", failureCount)
	if failureCount > 0 {
		os.Exit(1)
	}
}

func runOne(runner TestRunner) (err error) {
	name := runner.Name()
	fmt.Printf("-------------- executing test: %s --------------\n", name)

	defer func() {
		if err = runner.Cleanup(); err != nil {
			fmt.Printf("%s::Cleanup() returned error: %s\n", name, err)
		}
	}()

	if err = runner.Setup(); err != nil {
		fmt.Printf("%s::Setup() returned error: %s\n", name, err)
	}
	if err = runner.Run(); err != nil {
		fmt.Printf("%s::Run() returned error: %s\n", name, err)
	}
	return
}
