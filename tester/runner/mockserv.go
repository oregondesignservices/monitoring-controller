package main

import (
	"bytes"
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/oregondesignservices/monitoring-controller/tester/common"
	"net/http"
	"time"
)

const mockServerUrl = "http://localhost:9091"

func responseIsOk(r *http.Response) bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

func ResetMockServerData() error {
	response, err := httpClient.Post(mockServerUrl+"/reset", "", nil)
	if err != nil {
		return err
	}
	if !responseIsOk(response) {
		return errors.New("failed to reset mock server data")
	}
	return nil
}

func SetMockServerResponse(mock common.MockResponse) error {
	data, err := jsoniter.Marshal(mock)
	if err != nil {
		return err
	}

	response, err := httpClient.Post(mockServerUrl+"/set", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	if !responseIsOk(response) {
		return errors.New("failed to set mock server data")
	}
	return nil
}

func GetMockServerCapturedRequests() (common.CapturedRequests, error) {
	response, err := httpClient.Get(mockServerUrl + "/data")
	if err != nil {
		return nil, err
	}
	if !responseIsOk(response) {
		return nil, errors.New("failed to get mock server data")
	}
	defer response.Body.Close()

	captured := common.CapturedRequests{}
	return captured, jsoniter.NewDecoder(response.Body).Decode(&captured)
}

func WaitForAnyRequest(ctx context.Context) (common.CapturedRequests, error) {
	done := ctx.Done()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			data, err := GetMockServerCapturedRequests()
			if err != nil {
				return nil, err
			}
			if len(data) > 0 {
				return data, nil
			}
		case <-done:
			return nil, errors.New("no requests hit mock server in time")
		}
	}
}
