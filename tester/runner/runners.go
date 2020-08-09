package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/oregondesignservices/monitoring-controller/tester/common"
	"time"
)

type TestRunner interface {
	Name() string
	Setup() error
	Run() error
	Cleanup() error
}

// Does nothing but fill in the TestRunner interface so tests do not have to write all
// 3 functions unless they have to
type NoopRunner struct{}

func (n NoopRunner) Setup() error {
	return nil
}

func (n NoopRunner) Run() error {
	return nil
}

func (n NoopRunner) Cleanup() error {
	return nil
}

type TestSimpleApply struct {
	NoopRunner
	PathToApply string
}

func (n *TestSimpleApply) Name() string {
	return fmt.Sprintf("TestSimpleApply{%s}", n.PathToApply)
}

func (n *TestSimpleApply) Setup() error {
	return KubectlApply(n.PathToApply)
}

func (n *TestSimpleApply) Cleanup() error {
	return KubectlDelete(n.PathToApply)
}

type WaitForController struct {
	NoopRunner
	HealthUrl string
}

func (n *WaitForController) Name() string {
	return fmt.Sprintf("WaitForController{%s}", n.HealthUrl)
}

func (n *WaitForController) Run() error {
	for i := 0; i < 12; i++ {
		resp, err := httpClient.Get(n.HealthUrl)
		if err == nil && resp.StatusCode == 200 {
			logger.Infof("%s is awake!", n.HealthUrl)
			return nil
		} else {
			logger.Infof("waiting for %s to awaken...", n.HealthUrl)
		}
		time.Sleep(5 * time.Second)
	}
	return errors.New("server never woke up")
}

type TestSingleHttpMonitorEndpoint struct {
	*TestSimpleApply
	MockedResponse  common.MockResponse
	ExpectedRequest common.CapturedRequest
}

func (t *TestSingleHttpMonitorEndpoint) Name() string {
	return fmt.Sprintf("TestSingleHttpMonitorEndpoint{%s, %s}", t.MockedResponse.Method, t.MockedResponse.UriPath)
}

func (t *TestSingleHttpMonitorEndpoint) Setup() error {
	if err := t.TestSimpleApply.Setup(); err != nil {
		return err
	}
	// setup the response mock
	return SetMockServerResponse(t.MockedResponse)
}

func (t *TestSingleHttpMonitorEndpoint) Run() error {
	// wait for the HttpMonitor to run
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	captured, err := WaitForAnyRequest(ctx)
	if err != nil {
		return err
	}

	single, exists := captured.Get(t.MockedResponse.UriPath, t.MockedResponse.Method)
	if !exists {
		return fmt.Errorf("%s::%s was never called", t.MockedResponse.UriPath, t.MockedResponse.Method)
	}

	if !IsSubset(t.ExpectedRequest.LastQueryParams, single.LastQueryParams) {
		return errors.New("LastQueryParams not expected")
	}
	if !IsSubset(t.ExpectedRequest.LastRequestHeaders, single.LastRequestHeaders) {
		return errors.New("LastQueryParams not expected")
	}

	return nil
}

func (t *TestSingleHttpMonitorEndpoint) Cleanup() error {
	if err := t.TestSimpleApply.Cleanup(); err != nil {
		return err
	}
	// clear response data so the next test is fresh
	return ResetMockServerData()
}
