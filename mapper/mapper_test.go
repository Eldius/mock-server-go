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
