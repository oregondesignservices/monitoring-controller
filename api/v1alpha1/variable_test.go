package v1alpha1

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func newReaderCloser(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

func TestVariable_ParseFromResponse(t *testing.T) {
	tests := []struct {
		TestName string
		Var *Variable
		Resp *http.Response
		ExpectErr bool
		ExpectedValue string
	}{
		// Header tests
		{
			"get-header-simple",
			&Variable{
				Name: "test",
				From: FromTypeHeaders,
				JsonPath: "/My-Header",
			},
			&http.Response{
				Header: http.Header{
					"My-Header": []string{"val"},
				},
			},
			false,
			"val",
		},
		{
			"get-header-by-index",
			&Variable{
				Name: "test",
				From: FromTypeHeaders,
				JsonPath: "/My-Header/0",
			},
			&http.Response{
				Header: http.Header{
					"My-Header": []string{"val"},
				},
			},
			false,
			"val",
		},
		{
			"get-header-by-index-1",
			&Variable{
				Name: "test",
				From: FromTypeHeaders,
				JsonPath: "/My-Header/1",
			},
			&http.Response{
				Header: http.Header{
					"My-Header": []string{"val", "val2"},
				},
			},
			false,
			"val2",
		},
		{
			"header-404",
			&Variable{
				Name: "test",
				From: FromTypeHeaders,
				JsonPath: "/Not-Real",
			},
			&http.Response{
				Header: http.Header{
					"My-Header": []string{"val", "val2"},
				},
			},
			true,
			"",
		},
		// JSON tests
		{
			"json-simple",
			&Variable{
				Name: "test",
				From: FromTypeBodyJson,
				JsonPath: "/myVar",
			},
			&http.Response{
				Body: newReaderCloser(`{"myVar": "abc"}`),
			},
			false,
			"abc",
		},
		{
			"json-array",
			&Variable{
				Name: "test",
				From: FromTypeBodyJson,
				JsonPath: "/myVar/0",
			},
			&http.Response{
				Body: newReaderCloser(`{"myVar": ["xyz"]}`),
			},
			false,
			"xyz",
		},
		{
			"json-array-1",
			&Variable{
				Name: "test",
				From: FromTypeBodyJson,
				JsonPath: "/myVar/1",
			},
			&http.Response{
				Body: newReaderCloser(`{"myVar": ["abc", "0"]}`),
			},
			false,
			"0",
		},
		{
			"json-int-but-is-string",
			&Variable{
				Name: "test",
				From: FromTypeBodyJson,
				JsonPath: "/intval",
			},
			&http.Response{
				Body: newReaderCloser(`{"intval": 12}`),
			},
			false,
			"12",
		},
		{
			"json-bool-but-is-string",
			&Variable{
				Name: "test",
				From: FromTypeBodyJson,
				JsonPath: "/boolval",
			},
			&http.Response{
				Body: newReaderCloser(`{"boolval": true}`),
			},
			false,
			"true",
		},
		{
			"json-404",
			&Variable{
				From: FromTypeBodyJson,
				JsonPath: "/notreal",
			},
			&http.Response{
				Body: newReaderCloser(`{"boolval": true}`),
			},
			true,
			"",
		},
		{
			"json-invalid-json",
			&Variable{
				From: FromTypeBodyJson,
				JsonPath: "/notreal",
			},
			&http.Response{
				Body: newReaderCloser(`junk`),
			},
			true,
			"",
		},
		// YAML test
		{
			"yaml-simple",
			&Variable{
				From: FromTypeBodyYaml,
				JsonPath: "/myVal",
			},
			&http.Response{
				Body: newReaderCloser(`
someVal: "abc"
myVal: "qwerty"
`),
			},
			false,
			"qwerty",
		},
		{
			"yaml-array",
			&Variable{
				Name: "test",
				From: FromTypeBodyYaml,
				JsonPath: "/myVal/1",
			},
			&http.Response{
				Body: newReaderCloser(`
someVal: "abc"
myVal: ["qwerty", "asdf"]
`),
			},
			false,
			"asdf",
		},
		{
			"yaml-404",
			&Variable{
				From: FromTypeBodyYaml,
				JsonPath: "/notreal/1",
			},
			&http.Response{
				Body: newReaderCloser(`
someVal: "abc"
myVal: ["qwerty", "asdf"]
`),
			},
			true,
			"",
		},
		{
			"yaml-invalid-yaml",
			&Variable{
				From: FromTypeBodyYaml,
				JsonPath: "/notreal/1",
			},
			&http.Response{
				Body: newReaderCloser(`junk`),
			},
			true,
			"",
		},
	}

	for _, testdata := range tests {
		err := testdata.Var.ParseFromResponse(testdata.Resp)
		if err == nil && testdata.ExpectErr {
			t.Errorf("[%s] expected error but got none", testdata.TestName)
			continue
		}
		if err != nil && !testdata.ExpectErr {
			t.Errorf("[%s] got unexpected err: %s", testdata.TestName, err)
			continue
		}

		if testdata.Var.Value != testdata.ExpectedValue {
			t.Errorf("[%s] value unexpected. Got: %s, expected: %s", testdata.TestName, testdata.Var.Value, testdata.ExpectedValue)
		}
	}
}
