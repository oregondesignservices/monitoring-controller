package main

import (
	"errors"
	"fmt"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
	return KubectlApply("default", n.PathToApply)
}

func (n *TestSimpleApply) Cleanup() error {
	return KubectlDelete("default", n.PathToApply)
}

type WaitForController struct {
	NoopRunner
}

func (n *WaitForController) Name() string {
	return "WaitForController{}"
}

func (n *WaitForController) Run() error {
	for i := 0; i < 12; i++ {
		resp, err := httpClient.Get(metricsUrl)
		if err == nil && resp.StatusCode == 200 {
			fmt.Println("controller is awake!")
			return nil
		} else {
			fmt.Println("waiting for controller to awaken...")
		}
		time.Sleep(5 * time.Second)
	}
	return errors.New("controller never woke up")
}

type EnsureSecretIsCreated struct {
	*TestSimpleApply
	SecretNamespace string
	SecretName      string
	Key             string
}

func (n *EnsureSecretIsCreated) Name() string {
	return fmt.Sprintf("EnsureSecretIsCreated{%s, %s, %s, %s}", n.PathToApply, n.SecretNamespace, n.SecretName, n.Key)
}

func (n *EnsureSecretIsCreated) Run() error {
	for i := 0; i < 10; i++ {
		secret, err := GetSecret(n.SecretNamespace, n.SecretName)
		if err != nil {
			if apierrors.IsNotFound(err) {
				fmt.Println("waiting for secret to be created...")
			} else {
				return err
			}
		} else {
			_, exists := secret.Data[n.Key]
			if exists {
				return nil
			}
			return errors.New("key not found in secret")
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("secret was never created")
}

type RejectInvalidCrd struct {
	NoopRunner
	InvalidPathToApply string
}

func (n *RejectInvalidCrd) Name() string {
	return fmt.Sprintf("RejectInvalidCrd{%s}", n.InvalidPathToApply)
}

func (n *RejectInvalidCrd) Run() error {
	err := KubectlApply("default", n.InvalidPathToApply)
	if err == nil {
		return errors.New("an invalid crd was applied without error")
	}
	return nil
}

func (n *RejectInvalidCrd) Cleanup() error {
	// We ignore the error because a valid test will have nothing to delete.
	// It's only a failed test that will need cleanup
	_ = KubectlDelete("default", n.InvalidPathToApply)
	return nil
}
