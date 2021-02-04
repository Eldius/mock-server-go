package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Eldius/mock-server-go/mapper"
)

var r mapper.Router
var c http.Client

func init() {
	r = mapper.ImportMappingYaml("../mapper/samples/mapping_file_benchmark.yml")
	c = http.Client{}
}

type requestBuilder func() (*http.Request, error)

func BenchmarkMapperResolve01(b *testing.B) {
	b.StopTimer()
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.Handle)

	server := httptest.NewServer(mux)
	defer server.Close()

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		switch n % 4 {
		case 0:
			testRequest(
				func() (*http.Request, error) {
					return http.NewRequest("GET", fmt.Sprintf("%s/v1/contract", server.URL), nil)
				},
				`{"id": 123, "name": "My Contract"}`,
				200,
				b,
			)
		case 1:
			testRequest(
				func() (*http.Request, error) {
					return http.NewRequest("GET", fmt.Sprintf("%s/v1/contract/not/found", server.URL), nil)
				},
				`Mapping not found`,
				404,
				b,
			)
		case 2:
			testRequest(
				func() (*http.Request, error) {
					return http.NewRequest("POST", fmt.Sprintf("%s/v1/test", server.URL), newReqBody(`{"key": "value"}`))
				},
				`{"id": 123, "name": "My Contract"}`,
				200,
				b,
			)
		case 3:
			testRequest(
				func() (*http.Request, error) {
					return http.NewRequest("POST", fmt.Sprintf("%s/v1/contract", server.URL), newReqBody(`{"key": "value"}`))
				},
				``,
				202,
				b,
			)

		}
	}
}

func newReqBody(body string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(body)))
}

func testRequest(reqFun requestBuilder, msg string, statusCode int, b *testing.B) {

	req, err := reqFun()
	if err != nil {
		b.Errorf("Failed to create request\n%s", err.Error())
	}
	res, err := c.Do(req)
	if err != nil {
		b.Errorf("Failed to execute request\n%s", err.Error())
	}
	defer res.Body.Close()

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		b.Errorf("Failed to read response body\n%s", err.Error())
	}
	resBody := string(byteBody)
	if res.StatusCode != statusCode {
		b.Fatalf(`Get response must return status code '%d', but was '%d'`, statusCode, res.StatusCode)
	}
	if resBody != msg {
		b.Fatalf(`Get response must return '%s', but was '%s'`, msg, resBody)
	}
}
