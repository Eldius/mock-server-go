package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Eldius/mock-server-go/mapper"
	"gopkg.in/yaml.v3"
)

const (
	mappingFile = "../mapper/samples/mapping_file_test.yml"
)

func TestRouteHandler(t *testing.T) {
	r := mapper.ImportMappingYaml(mappingFile)
	h := RouteHandler(&r)
	mux := http.NewServeMux()
	mux.HandleFunc("/route", h)

	server := httptest.NewServer(mux)
	defer server.Close()

	url := fmt.Sprintf("%s/route", server.URL)

	newRouteReq := strings.NewReader(`{
    "path": "/app/v1/test",
    "method": "GET",
    "response": {
        "headers": {},
        "body": "{\"msg\": \"Test response!\"}",
        "statusCode": 202
    }
}`)

	c := http.Client{}
	res, err := c.Post(url, "application/json", newRouteReq)
	if err != nil {
		t.Errorf("Failed to make request\n%s", err.Error())
	}
	defer res.Body.Close()

	resBodyBin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to parse response body\n%s", err.Error())
	}
	resBody := string(resBodyBin)
	t.Log(resBody)

	if res.StatusCode != 200 {
		t.Errorf("Status code for new route must be '200', but was '%d'", res.StatusCode)
	}
	if resBody == "" {
		t.Fatalf("Response body must not be nil, but was '%s'", resBody)
	}
	var mapping mapper.RequestMapping
	err = json.Unmarshal(resBodyBin, &mapping)
	if err != nil {
		t.Errorf("Failed to parse response body to object\n%s", err.Error())
	}

	if mapping.Method != "GET" {
		t.Errorf("Response mapping must have method 'GET', but was '%s'", mapping.Method)
	}
	if mapping.Path != "/app/v1/test" {
		t.Errorf("Response mapping must have method '/app/v1/test', but was '%s'", mapping.Path)
	}

	if len(r.Routes) != 3 {
		t.Errorf("Routes count must be 3, but was %d", len(r.Routes))
	}
}

func TestRouteHandlerError(t *testing.T) {
	r := mapper.ImportMappingYaml(mappingFile)
	h := RouteHandler(&r)
	mux := http.NewServeMux()
	mux.HandleFunc("/route", h)

	server := httptest.NewServer(mux)
	defer server.Close()

	url := fmt.Sprintf("%s/route", server.URL)

	newRouteReq := strings.NewReader(``)

	c := http.Client{}
	res, err := c.Post(url, "application/json", newRouteReq)
	if err != nil {
		t.Errorf("Failed to make request\n%s", err.Error())
	}
	defer res.Body.Close()

	resBodyBin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to parse response body\n%s", err.Error())
	}
	resBody := string(resBodyBin)
	t.Log(resBody)

	if res.StatusCode != 500 {
		t.Errorf("Status code for new route must be '200', but was '%d'", res.StatusCode)
	}
	if resBody == "" {
		t.Fatalf("Response body must not be nil, but was '%s'", resBody)
	}
	if len(r.Routes) != 2 {
		t.Errorf("Routes count must be 2, but was %d", len(r.Routes))
	}
}

func TestRouteHandlerGet(t *testing.T) {
	r := mapper.ImportMappingYaml(mappingFile)
	h := RouteHandler(&r)
	mux := http.NewServeMux()
	mux.HandleFunc("/route", h)

	server := httptest.NewServer(mux)
	defer server.Close()

	url := fmt.Sprintf("%s/route", server.URL)

	c := http.Client{}
	res, err := c.Get(url)
	if err != nil {
		t.Errorf("Failed to make request\n%s", err.Error())
	}
	defer res.Body.Close()

	var response []mapper.RequestMapping
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to unmarshall request body\n%s", err.Error())
	}

	if len(response) != 2 {
		t.Errorf("Must return '2' mapping objects, but returned '%d'", len(response))
	}

	if strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") {
		t.Errorf("Must return '2' mapping objects, but returned '%d'", len(response))
	}
}

func TestRouteHandlerGetYAML(t *testing.T) {
	r := mapper.ImportMappingYaml(mappingFile)
	h := RouteHandler(&r)
	mux := http.NewServeMux()
	mux.HandleFunc("/route", h)

	server := httptest.NewServer(mux)
	defer server.Close()

	url := fmt.Sprintf("%s/route", server.URL)

	c := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("Failed to create request object\n%s", err.Error())
	}
	req.Header.Add("Accept", "application/yaml")
	res, err := c.Do(req)
	if err != nil {
		t.Errorf("Failed to make request\n%s", err.Error())
	}
	defer res.Body.Close()

	var response []mapper.RequestMapping
	err = yaml.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to unmarshall request body\n%s", err.Error())
	}

	if len(response) != 2 {
		t.Errorf("Must return '2' mapping objects, but returned '%d'", len(response))
	}

	if strings.HasPrefix(res.Header.Get("Content-Type"), "application/yaml") {
		t.Errorf("Must return '2' mapping objects, but returned '%d'", len(response))
	}
}

func TestRouteHandlerGetJSON(t *testing.T) {
	r := mapper.ImportMappingYaml(mappingFile)
	h := RouteHandler(&r)
	mux := http.NewServeMux()
	mux.HandleFunc("/route", h)

	server := httptest.NewServer(mux)
	defer server.Close()

	url := fmt.Sprintf("%s/route", server.URL)

	c := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("Failed to create request object\n%s", err.Error())
	}
	req.Header.Add("Accept", "application/json")
	res, err := c.Do(req)
	if err != nil {
		t.Errorf("Failed to make request\n%s", err.Error())
	}
	if err != nil {
		t.Errorf("Failed to make request\n%s", err.Error())
	}
	defer res.Body.Close()

	var response []mapper.RequestMapping
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to unmarshall request body\n%s", err.Error())
	}

	if len(response) != 2 {
		t.Errorf("Must return '2' mapping objects, but returned '%d'", len(response))
	}

	if strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") {
		t.Errorf("Must return '2' mapping objects, but returned '%d'", len(response))
	}
}

func TestRouteHandlerMethodNotAllowed(t *testing.T) {
	r := mapper.ImportMappingYaml(mappingFile)
	h := RouteHandler(&r)
	mux := http.NewServeMux()
	mux.HandleFunc("/route", h)

	server := httptest.NewServer(mux)
	defer server.Close()

	url := fmt.Sprintf("%s/route", server.URL)

	c := http.Client{}
	res, err := c.Head(url)
	if err != nil {
		t.Errorf("Failed to make request\n%s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != 405 {
		t.Errorf("Response code must be '405', but was '%d'", res.StatusCode)
	}
}
