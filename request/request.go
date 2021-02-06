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
	Path    string  `json:"path"`
	Method  string  `json:"method"`
	Headers Headers `json:"headers"`
	Body    string  `json:"body"`
}

type ResponseRecord struct {
	Headers Headers `json:"headers"`
	Code    int     `json:"code"`
	Body    string  `json:"body"`
}

type Record struct {
	ID       int            `json:"id"`
	ReqID    uuid.UUID      `json:"reqId"`
	Request  RequestRecord  `json:"request"`
	Response ResponseRecord `json:"response"`
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
