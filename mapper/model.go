package mapper

import (
	"net/http"
	"strings"

	"github.com/Eldius/mock-server-go/request"
	//lua "github.com/yuin/gopher-lua"
)

type MockHeader map[string][]string

type MockResponse struct {
	Headers    MockHeader `json:"headers"`
	Body       *string    `json:"body"`
	Script     *string    `json:"script"`
	StatusCode int        `json:"statusCode"`
}

type RequestMapping struct {
	Path     string       `json:"path"`
	Method   string       `json:"method"`
	Response MockResponse `json:"response"`
}

func (r *MockResponse) IsScript() bool {
	return r.Script != nil
}

func (r *RequestMapping) MakeResponse(rw http.ResponseWriter, req *http.Request) request.ResponseRecord {
	respRec := request.ResponseRecord{
		Headers: map[string][]string{},
	}
	for k, values := range r.Response.Headers {
		respRec.Headers[k] = append(respRec.Headers[k], values...)
		rw.Header().Add(k, strings.Join(respRec.Headers[k], "; "))
	}
	if r.Response.IsScript() {
		resBody, resCode, err := r.parseScript(rw, req)
		if err != nil {
			respRec.Body = err.Error()
			respRec.Code = http.StatusInternalServerError
			log.WithError(err).Error("Failed to execute script")
		} else {
			respRec.Body = resBody
			respRec.Code = resCode
		}
		rw.WriteHeader(respRec.Code)
		_, _ = rw.Write([]byte(respRec.Body))
	} else {
		respRec.Body = *r.Response.Body
		rw.WriteHeader(r.Response.StatusCode)
		_, _ = rw.Write([]byte(*r.Response.Body))
	}
	return respRec
}
