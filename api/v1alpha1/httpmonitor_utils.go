package v1alpha1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func replaceQueryParams(v map[string][]string, replacer *strings.Replacer) url.Values {
	if len(v) == 0 {
		return v
	}
	newValues := make(url.Values)

	for key, values := range v {
		for _, v := range values {
			newValues.Add(key, replacer.Replace(v))
		}
	}

	return newValues
}

func replaceHeader(v http.Header, replacer *strings.Replacer) http.Header {
	if len(v) == 0 {
		return v
	}

	newHeaders := make(http.Header)

	for key, values := range v {
		for _, v := range values {
			newHeaders.Add(key, replacer.Replace(v))
		}
	}

	return newHeaders
}

func (r *HttpRequest) BuildRequest() (*http.Request, error) {
	replacer := r.availableVariables.newReplacer()

	finalUrl := replacer.Replace(r.Url)
	body := replacer.Replace(r.Body)
	query := replaceQueryParams(r.QueryParams, replacer)
	header := replaceHeader(r.Headers, replacer)

	req, err := http.NewRequest(r.Method, finalUrl, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header = header

	req.URL.RawQuery = query.Encode()
	return req, nil
}

func containsInt(needle int, haystay []int) bool {
	for _, val := range haystay {
		if val == needle {
			return true
		}
	}
	return false
}

// Send the HTTP request
func (r *HttpRequest) Do(client *http.Client) (*http.Response, error) {
	req, err := r.BuildRequest()
	if err != nil {
		return nil, err
	}

	timeoutDuration, err  := time.ParseDuration(r.Timeout)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	return client.Do(req.WithContext(ctx))
}

func (r *HttpRequest) HandleResponse(resp *http.Response) error {
	if !containsInt(resp.StatusCode, r.ExpectedResponseCodes) {
		return fmt.Errorf("not an expected error code: %d is not in %x", resp.StatusCode, r.ExpectedResponseCodes)
	}
	// Nothing to parse
	if len(r.Variables) == 0 {
		return nil
	}


	return nil
}
