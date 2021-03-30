package mapper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Eldius/mock-server-go/config"
)

var router Router

func init() {
	mappingBody1 := `{
		"msg": "Success!"
	}`
	mappingBody2 := `{
		"msg": "Success!"
	}`
	var mapping []RequestMapping
	mapping = append(mapping, RequestMapping{
		Path:   "/path/test",
		Method: "POST",
		Response: MockResponse{
			StatusCode: 200,
			Body:       &mappingBody1,
			Headers: MockHeader{
				"Content-Type": "application/json",
			},
		},
	})
	router = NewRouter(mapping)
	router.Add(RequestMapping{
		Path:   "/path/test/get",
		Method: "GET",
		Response: MockResponse{
			StatusCode: 200,
			Body:       &mappingBody2,
			Headers: MockHeader{
				"Content-Type": "application/json",
			},
		},
	})
	config.Setup("")
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

	if len(r.Routes) != 3 {
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
	if *r.Routes[1].Response.Body != `{"id": 123, "name": "My Contract"}` {
		t.Errorf(`First mapping must have body '{"id": 123, "name": "My Contract"}', but has '%s'`, *r.Routes[1].Response.Body)
	}

	if r.Routes[2].Method != "POST" {
		t.Errorf("First mapping must be a 'POST', but has '%s'", r.Routes[2].Method)
	}
	if r.Routes[2].Path != "/v2/test" {
		t.Errorf("First mapping must have path '/v1/contract', but has '%s'", r.Routes[2].Path)
	}
	if !strings.Contains(*r.Routes[2].Response.Script, `res.body = JSON.stringify({"PING": "pong"});`) {
		t.Errorf(`Scripted mapping body must start with 'script:javascript:' prefix, but hasn't ('%s')`, *r.Routes[2].Response.Body)
	}
}

func TestMakeResponse(t *testing.T) {
	sourcePath := "samples/mapping_file_test.yml"
	r := ImportMappingYaml(sourcePath)
	rec := httptest.NewRecorder()
	r.Routes[0].MakeResponse(rec, httptest.NewRequest("POST", "/v1/contract", nil))
	respBody := rec.Body.String()
	if respBody != "" {
		t.Errorf("Response body must be empty, but was '%s'", respBody)
	}
	if rec.Code != 202 {
		t.Errorf("Response code must be '202', but was '%d'", rec.Code)
	}

	rec1 := httptest.NewRecorder()
	r.Routes[1].MakeResponse(rec1, httptest.NewRequest("GET", "/v1/contract", nil))
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

func TestRequestMappingParseScript(t *testing.T) {
	script := `
	console.log(req.body);
	var body = JSON.parse(req.body);
	console.log(body);
	var res = {};
	if (body.contract) {
	  res.code = 200;
	  res.body = JSON.stringify({
		"contract": body.contract,
		"status": "OK"
	  });
	} else {
	  res.code = 200;
	  res.body = JSON.stringify({"PING": "pong"});
	}`
	r := RequestMapping{
		Method: "POST",
		Response: MockResponse{
			Script: &script,
		},
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/test", bytes.NewBuffer([]byte(`{"contract": 1234}`)))
	body, code, err := r.parseScript(rw, req)
	if err != nil {
		t.Errorf("Failed to execute script\n%s\n", err)
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal([]byte(body), &bodyMap)
	if err != nil {
		t.Errorf("Failed to parse body to a map\n%s\n", err)
	}
	if bodyMap["contract"] != float64(1234) {
		t.Errorf("body.contract should be '1234', but was '%d'", bodyMap["contract"])
	}
	if bodyMap["status"] != "OK" {
		t.Errorf("body.contract should be 'OK', but was '%s'", bodyMap["status"])
	}

	if code != 200 {
		t.Errorf("Response code should be '200', but was '%d'", code)
	}
}

func TestRequestMappingParseScriptExecutionError(t *testing.T) {
	script := `script:javascript:
	throw new Error("Holy crap!")
`
	r := RequestMapping{
		Method: "POST",
		Response: MockResponse{
			Script: &script,
		},
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/test", bytes.NewBufferString(`{"contract": 1234}`))
	_, _, err := r.parseScript(rw, req)
	if err == nil {
		t.Errorf("Failed to execute script\n%s\n", err)
	}

	if !strings.Contains(err.Error(), `throw new Error("Holy crap!")`) {
		t.Error("Exception must contains script content")
	}
}

func TestRequestMappingParseScriptInvalidScript(t *testing.T) {
	script := `script:javascript:
	_1234;
`
	r := RequestMapping{
		Method: "POST",
		Response: MockResponse{
			Script: &script,
		},
	}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/test", bytes.NewBuffer([]byte(`{"contract": 1234}`)))
	_, _, err := r.parseScript(rw, req)
	if err == nil {
		t.Errorf("Failed to execute script\n%s\n", err)
	}
	if !strings.Contains(err.Error(), `_1234;`) {
		t.Error("Exception must contains script content")
	}

}

func TestHandleRequestScripted(t *testing.T) {
	r := ImportMappingYaml(mappingFile)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.Handle)

	server := httptest.NewServer(mux)
	defer server.Close()

	url := fmt.Sprintf("%s/v2/test", server.URL)

	c := http.Client{}
	res, err := c.Post(url, "application/json", bytes.NewBufferString(`{
		"contract": 12345,
		"status": "pending"
	}`))
	if err != nil {
		t.Fatalf("Failed to execute request\n%s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf(`Scripted response must return status code '200', but was '%d'`, res.StatusCode)
	}
	var bodyMap map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&bodyMap)
	if err != nil {
		t.Errorf("Failed to parse body to a map\n%s\n", err)
	}
	if bodyMap["contract"] != float64(12345) {
		t.Fatalf(`Returned contract must be equals to '12345' , but was '%d'`, bodyMap["contract"])
	}
	resContentType := res.Header.Get("Content-Type")
	if !strings.HasPrefix(resContentType, "application/json") {
		t.Fatalf(`Get response must contain header 'Content-Type' equals to 'application/json', but has '%s'`, resContentType)
	}
}
