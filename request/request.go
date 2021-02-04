package request

import (
	"bytes"
	"io/ioutil"
	"log"
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
	Headers Headers
	Code    int
	Body    string
}

type Record struct {
	ID       int
	ReqID    uuid.UUID
	Request  RequestRecord
	Response ResponseRecord
}

func NewRecord(r *http.Request) *Record {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body\n%s", err.Error())
		return nil
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	log.Printf("---\nrequest body:\n%s\n---", string(body))

	return &Record{
		ReqID: uuid.New(),
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
