package mapper

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var router Router

func init() {
	var mapping []RequestMapping
	mapping = append(mapping, RequestMapping{
		Path:   "/path/test",
		Method: "POST",
		Response: MockResponse{
			StatusCode: 200,
			Body: `{
					"msg": "Success!"
				}`,
			Headers: MockHeader{
				"Content-Type": []string{"application/json"},
			},
		},
	})
	router = NewRouter(mapping)
	router.Add(RequestMapping{
		Path:   "/path/test/get",
		Method: "GET",
		Response: MockResponse{
			StatusCode: 200,
			Body: `{
					"msg": "Success!"
				}`,
			Headers: MockHeader{
				"Content-Type": []string{"application/json"},
			},
		},
	})
}

func TestRouteFound(t *testing.T) {

	reqBody := bufio.NewReader(strings.NewReader(`{"key": "value"}`))
	req := httptest.NewRequest("POST", "http://localhost:8081/path/test", reqBody)

	mapping := router.Route(req)

	if mapping == nil {
		t.Errorf("Mapping must not be nil")
	}
}

func TestRouteFound1(t *testing.T) {

	reqBody := bufio.NewReader(strings.NewReader(`{"key": "value"}`))
	req := httptest.NewRequest("GET", "http://localhost:8081/path/test/get", reqBody)

	mapping := router.Route(req)

	if mapping == nil {
		t.Errorf("Mapping must not be nil")
	}
}

func TestRouteFoundUpperCase(t *testing.T) {

	reqBody := bufio.NewReader(strings.NewReader(`{"key": "value"}`))
	req := httptest.NewRequest("POST", "http://localhost:8081/PATH/TEST", reqBody)

	mapping := router.Route(req)

	if mapping == nil {
		t.Errorf("Mapping must not be nil")
	}
}

func TestRouteNotFound(t *testing.T) {

	reqBody := bufio.NewReader(strings.NewReader(`{"key": "value"}`))
	req := httptest.NewRequest("POST", "http://localhost:8081/path/not/exists", reqBody)

	mapping := router.Route(req)

	if mapping != nil {
		t.Errorf("Mapping must be nil")
	}
}

func TestImportMappingYaml(t *testing.T) {
	sourcePath := "samples/mapping_file_test.yml"
	r := ImportMappingYaml(sourcePath)

	if len(r.Routes) != 2 {
		t.Errorf("Must have 2 mappings, but has '%d'", len(r.Routes))
	}

	if r.Routes[0].Method != "POST" {
		t.Errorf("First mapping must be a 'POST', but has '%s'", r.Routes[0].Method)
	}
	if r.Routes[0].Path != "/v1/contract" {
		t.Errorf("First mapping must have path '/v1/contract', but has '%s'", r.Routes[0].Path)
	}

	if r.Routes[1].Method != "GET" {
		t.Errorf("First mapping must be a 'GET', but has '%s'", r.Routes[1].Method)
	}
	if r.Routes[1].Path != "/v1/contract" {
		t.Errorf("First mapping must have path '/v1/contract', but has '%s'", r.Routes[1].Path)
	}
	if r.Routes[1].Response.Body != `{"id": 123, "name": "My Contract"}` {
		t.Errorf(`First mapping must have body '{"id": 123, "name": "My Contract"}', but has '%s'`, r.Routes[1].Response.Body)
	}
}

func TestMakeResponse(t *testing.T) {
	sourcePath := "samples/mapping_file_test.yml"
	r := ImportMappingYaml(sourcePath)
	rec := httptest.NewRecorder()
	r.Routes[0].MakeResponse(rec)
	respBody := rec.Body.String()
	if respBody != "" {
		t.Errorf("Response body must be empty, but was '%s'", respBody)
	}
	if rec.Code != 202 {
		t.Errorf("Response code must be '202', but was '%d'", rec.Code)
	}

	rec1 := httptest.NewRecorder()
	r.Routes[1].MakeResponse(rec1)
	respBody1 := rec1.Body.String()
	if respBody1 != `{"id": 123, "name": "My Contract"}` {
		t.Errorf(`Response body must be '{"id": 123, "name": "My Contract"}', but was '%s'`, respBody1)
	}
	if rec1.Code != 200 {
		t.Errorf("Response code must be '200', but was '%d'", rec1.Code)
	}
}

const (
	mappingFile = "../mapper/samples/mapping_file_test.yml"
)

func TestHandleRequestGet(t *testing.T) {
	r := ImportMappingYaml(mappingFile)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.Handle)

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
	r := ImportMappingYaml(mappingFile)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.Handle)

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
	r := ImportMappingYaml(mappingFile)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.Handle)

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
