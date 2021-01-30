package mapper

import "net/http"

type MockHeader map[string][]string

type MockResponse struct {
	Headers    MockHeader `json:"headers"`
	Body       string     `json:"body"`
	StatusCode int        `json:"statusCode"`
}

type RequestMapping struct {
	Path     string       `json:"path"`
	Method   string       `json:"method"`
	Response MockResponse `json:"response"`
}

func (r *RequestMapping) MakeResponse(rw http.ResponseWriter) {
	for k, values := range r.Response.Headers {
		for _, v := range values {
			rw.Header().Add(k, v)
		}
	}

	rw.WriteHeader(r.Response.StatusCode)
	_, _ = rw.Write([]byte(r.Response.Body))
}
