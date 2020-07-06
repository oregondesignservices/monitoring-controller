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
	"bytes"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (v VariableList) newReplacer() *strings.Replacer {
	args := make([]string, 2*len(v))
	for i := range v {
		args[0] = "{" + v[i].Name + "}"
		args[1] = v[i].Value
	}
	return strings.NewReplacer(args...)
}

// Clears all values that are not provided by users
func (v VariableList) clearValues() {
	for _, elem := range v {
		if elem.From != FromTypeProvided {
			elem.Value = ""
		}
	}
}

// Read a Response.Body and reset it for later reading.
// Required if we need to read a response body more than once.
func readBodyAndReset(resp *http.Response) []byte {
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close() //  must close, or we might have a memory leak
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

func (v *Variable) ParseFromResponse(resp *http.Response) error {
	switch v.From {
	case FromTypeProvided:
		// "provided" means the value is provided by the user
		return nil
	case FromTypeBodyJson:
		return v.parseFromBodyJson(resp)
	case FromTypeBodyYaml:
		return v.parseFromBodyYaml(resp)
	case FromTypeBodyRaw:
		return v.parseFromBodyRaw(resp)
	case FromTypeHeaders:
		return v.parseFromHeaders(resp)
	}
	return fmt.Errorf("not a known variable 'from' type: %s", v.From)
}

func (v *Variable) jsonPathToPieces() []string {
	pieces := strings.Split(v.JsonPath, "/")
	var finalPieces []string

	for _, piece := range pieces {
		if piece != "" {
			finalPieces = append(finalPieces, piece)
		}
	}
	return finalPieces
}

func (v *Variable) parseFromJsonBytes(jsonBody []byte) error {
	jsonPath := v.jsonPathToPieces()

	// jsoniter.Get needs a specific type. So convert to that.
	interfaceJsonPath := make([]interface{}, len(jsonPath))
	for i, val := range jsonPath {
		intVal, err := strconv.Atoi(val)
		if err == nil {
			interfaceJsonPath[i] = intVal
		} else {
			interfaceJsonPath[i] = val
		}
	}

	getter := jsoniter.Get(jsonBody, interfaceJsonPath...)
	err := getter.LastError()
	if err != nil {
		return err
	}
	v.Value = getter.ToString()
	return nil
}

func (v *Variable) parseFromBodyJson(resp *http.Response) error {
	body := readBodyAndReset(resp)
	return v.parseFromJsonBytes(body)
}

func (v *Variable) parseFromBodyYaml(resp *http.Response) error {
	// We convert YAML to json and then use jsoniter to do the same parsing.
	yamlBody := readBodyAndReset(resp)

	jsonBody, err := yaml.YAMLToJSON(yamlBody)
	if err != nil {
		return err
	}
	return v.parseFromJsonBytes(jsonBody)
}

func (v *Variable) parseFromBodyRaw(resp *http.Response) error {
	body := readBodyAndReset(resp)
	v.Value = string(body)
	return nil
}

func (v *Variable) parseFromHeaders(resp *http.Response) error {
	jsonBuf := &bytes.Buffer{}
	err := jsoniter.NewEncoder(jsonBuf).Encode(resp.Header)
	if err != nil {
		return err
	}

	pieces := v.jsonPathToPieces()
	switch len(pieces) {
	case 0:
		return errors.New("no valid jsonpath provided")
	case 1:
		// Assume "/something" is for the first element in a header
		v.Value = resp.Header.Get(pieces[0])
	case 2:
		// Assume "/something/1" wants a specific index of a header
		index, err := strconv.Atoi(pieces[1])
		if err != nil {
			return err
		}
		values := resp.Header.Values(pieces[0])
		if len(values) >= index {
			v.Value = values[index]
		}
	default:
		return fmt.Errorf("cannot parse jsonpath for header variable: %s", v.JsonPath)
	}

	if v.Value == "" {
		return fmt.Errorf("not a known header jsonpath: %s", v.JsonPath)
	}

	return nil
}
