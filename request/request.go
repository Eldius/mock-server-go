package request

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

type Headers map[string][]string

type RequestRecord struct {
	Path    string
	Method  string
	Headers Headers
	Body    string
}

type ResponseRecord struct {
	Path    string
	Method  string
	Headers Headers
	Code    int
	Body    string
}

type Record struct {
	ID       uuid.UUID
	Request  RequestRecord
	Response ResponseRecord
}

func NewRecord(r *http.Request) *Record {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return &Record{
		ID: uuid.New(),
		Request: RequestRecord{
			Path:   r.URL.Path,
			Method: r.Method,
			Body:   string(body),
		},
	}
}

func (r *Record) AddResponse(response ResponseRecord) *Record {
	r.Response = response
	return r
}
