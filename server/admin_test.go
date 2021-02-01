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
