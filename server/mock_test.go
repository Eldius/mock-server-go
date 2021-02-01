package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Eldius/mock-server-go/mapper"
)

const (
	mappingFile = "../mapper/samples/mapping_file_test.yml"
)

func TestHandleRequestGet(t *testing.T) {
	r := mapper.ImportMappingYaml(mappingFile)
	h := HandleMockRequest(&r)
	mux := http.NewServeMux()
	mux.HandleFunc("/", h)

	server := httptest.NewServer(mux)
	defer server.Close()

	c := http.Client{}
	res, err := c.Get(fmt.Sprintf("%s/v1/contract", server.URL))
	if err != nil {
		t.Fatalf("Failed to execute request\n%s", err.Error())
	}
	defer res.Body.Close()

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to parse response body\n%s", err.Error())
	}
	resBody := string(byteBody)

	if res.StatusCode != 200 {
		t.Fatalf(`Get response must return status code '200', but was '%d'`, res.StatusCode)
	}
	if resBody != `{"id": 123, "name": "My Contract"}` {
		t.Fatalf(`Get response must return '{"id": 123, "name": "My Contract"}', but was '%s'`, resBody)
	}
	if res.Header.Get("Content-Type") != "application/json" {
		t.Fatalf(`Get response must contain header 'Content-Type: application/jsob', but was '%s'`, res.Header.Get("Content-Type"))
	}
}

func TestHandleRequestPost(t *testing.T) {
	r := mapper.ImportMappingYaml(mappingFile)
	h := HandleMockRequest(&r)
	mux := http.NewServeMux()
	mux.HandleFunc("/", h)

	server := httptest.NewServer(mux)
	defer server.Close()

	url := fmt.Sprintf("%s/v1/contract", server.URL)

	c := http.Client{}
	res, err := c.Post(url, "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("Failed to execute request\n%s", err.Error())
	}
	defer res.Body.Close()

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to make request\n%s", err.Error())
	}
	resBody := string(byteBody)

	if res.StatusCode != 202 {
		t.Fatalf(`Get response must return status code '202', but was '%d'`, res.StatusCode)
	}
	if resBody != `` {
		t.Fatalf(`Get response must be empty, but was '%s'`, resBody)
	}
	if res.Header.Get("Content-Type") != "" {
		t.Fatalf(`Get response must not contain header 'Content-Type', but has '%s'`, res.Header.Get("Content-Type"))
	}
}

func TestHandleRequestGotNotFound(t *testing.T) {
	r := mapper.ImportMappingYaml(mappingFile)
	h := HandleMockRequest(&r)
	mux := http.NewServeMux()
	mux.HandleFunc("/", h)

	server := httptest.NewServer(mux)
	defer server.Close()

	url := fmt.Sprintf("%s/app/v1/not-found", server.URL)

	c := http.Client{}
	res, err := c.Get(url)
	if err != nil {
		t.Fatalf("Failed to execute request\n%s", err.Error())
	}
	defer res.Body.Close()

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to make request\n%s", err.Error())
	}
	resBody := string(byteBody)

	if res.StatusCode != 404 {
		t.Fatalf(`Get response must return status code '202', but was '%d'`, res.StatusCode)
	}
	if resBody != `Mapping not found` {
		t.Fatalf(`Get response be 'Mapping not found', but was '%s'`, resBody)
	}
	resContentType := res.Header.Get("Content-Type")
	if !strings.HasPrefix(resContentType, "text/plain") {
		t.Fatalf(`Get response must contain header 'Content-Type' equals to 'text/plain', but has '%s'`, resContentType)
	}
}
