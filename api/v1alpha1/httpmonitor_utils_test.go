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
package v1alpha1

import (
	"net/http"
	"net/url"
	"testing"
)

func TestVariableList_newReplacer(t *testing.T) {
	tests := []struct {
		InputString    string
		Vars           VariableList
		ExpectedOutput string
	}{
		{
			"http://test.com/{v1}",
			VariableList{
				&Variable{
					Name:  "v1",
					Value: "val1",
				},
			},
			"http://test.com/val1",
		},
	}

	for i, testdata := range tests {
		replacer := testdata.Vars.newReplacer()
		out := replacer.Replace(testdata.InputString)
		if out != testdata.ExpectedOutput {
			t.Errorf("[%d] unexpected output. Got: '%s', expected: '%s'", i, out, testdata.ExpectedOutput)
		}
	}
}

func TestHttpRequest_BuildRequest(t *testing.T) {
	queryParam := url.Values{}
	queryParam.Add("k", "v")
	queryParam.Add("k", "v")
	queryParam.Add("v1", "{v1}")

	header := http.Header{}
	header.Add("h", "val")
	header.Add("h", "val")
	header.Add("v1", "{v1}")

	availableVariables := VariableList{
		&Variable{
			Name:  "v1",
			Value: "val1",
		},
	}

	r := &HttpRequest{
		Url:                "http://test.com/{v1}",
		QueryParams:        queryParam,
		Headers:            header,
		AvailableVariables: availableVariables,
	}

	req, err := r.BuildRequest()
	if err != nil {
		t.Errorf("got err while building request: %s", err)
		return
	}

	expectedUrl := "http://test.com/val1?k=v&k=v&v1=val1"
	if req.URL.String() != expectedUrl {
		t.Errorf("unexpected url. Got: %s, wanted: %s", req.URL.String(), expectedUrl)
	}
}
