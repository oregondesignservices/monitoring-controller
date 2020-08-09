package main

import (
	"github.com/oregondesignservices/monitoring-controller/tester/common"
	"go.uber.org/zap"
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
	logger        *zap.SugaredLogger

	testdataDir = filepath.Join("tester", "runner", "testdata")
	samplesDir  = filepath.Join("config", "samples")
)

const (
	metricsUrl = "http://localhost:9090/metrics"
	kubeconfig = "kind-kubeconfig.yaml"
)

func main() {
	var err error
	plainLogger, err := zap.NewDevelopment()
	if err != nil {
		panic("failed to build logger: " + err.Error())
	}
	logger = plainLogger.Sugar()

	kubeClientset, err = newKubeClientset()
	if err != nil {
		logger.Panicf("failed to construct kubernetes clientset: %s", err)
	}

	runners := []TestRunner{
		// Wait for controller to wake up so we start real tests
		&WaitForController{HealthUrl: metricsUrl},
		// Wait for the mock server to wake up
		&WaitForController{HealthUrl: mockServerUrl + "/health"},

		&TestSingleHttpMonitorEndpoint{
			TestSimpleApply: &TestSimpleApply{
				PathToApply: filepath.Join(testdataDir, "simplest-possible.yaml"),
			},
			MockedResponse: common.MockResponse{
				UriPath: "/simplest",
				Method:  "GET",
				Status:  http.StatusOK,
			},
			ExpectedRequest: common.CapturedRequest{},
		},
		&TestSingleHttpMonitorEndpoint{
			TestSimpleApply: &TestSimpleApply{
				PathToApply: filepath.Join(testdataDir, "simplest-possible-with-var.yaml"),
			},
			MockedResponse: common.MockResponse{
				UriPath: "/simplest-with-var",
				Method:  "GET",
				Status:  http.StatusOK,
			},
			ExpectedRequest: common.CapturedRequest{},
		},
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

	logger.Info("-------------------------------------------------")
	logger.Infof("successful tests: %d", successCount)
	logger.Infof("failed tests: %d", failureCount)
	if failureCount > 0 {
		os.Exit(1)
	}
}

func runOne(runner TestRunner) (err error) {
	name := runner.Name()
	logger.Infof("-------------- executing test: %s --------------", name)

	defer func() {
		err2 := runner.Cleanup()
		if err2 != nil {
			logger.Errorf("%s::Cleanup() returned error: %s", name, err)
			err = err2
		}
	}()

	err = runner.Setup()
	if err != nil {
		logger.Errorf("%s::Setup() returned error: %s", name, err)
		return err
	}
	err = runner.Run()
	if err != nil {
		logger.Errorf("%s::Run() returned error: %s", name, err)
	}
	return
}
