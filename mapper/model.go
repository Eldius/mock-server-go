package mapper

import (
	"net/http"

	"github.com/Eldius/mock-server-go/request"
)

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

func (r *RequestMapping) MakeResponse(rw http.ResponseWriter) request.ResponseRecord {
	respRec := request.ResponseRecord{
		Headers: map[string][]string{},
	}
	for k, values := range r.Response.Headers {
		for _, v := range values {
			respRec.Headers[k] = append(respRec.Headers[k], v)
			rw.Header().Add(k, v)
		}
	}

	respRec.Code = r.Response.StatusCode
	rw.WriteHeader(r.Response.StatusCode)
	respRec.Body = r.Response.Body
	_, _ = rw.Write([]byte(r.Response.Body))
	return respRec
}
