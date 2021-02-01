package mapper

import (
	"bufio"
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
