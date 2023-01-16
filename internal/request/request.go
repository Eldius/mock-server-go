package request

import (
    "bytes"
    "io/ioutil"
    "net/http"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
)

type Headers map[string]string

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
    ID          int            `json:"id"`
    RequestDate time.Time      `json:"requestDate"`
    ReqID       uuid.UUID      `json:"reqId"`
    Request     RequestRecord  `json:"request"`
    Response    ResponseRecord `json:"response"`
}

func NewRecord(r *http.Request) *Record {
    bodyReader := r.Body
    defer bodyReader.Close()
    body, err := ioutil.ReadAll(bodyReader)
    if err != nil {
        log.WithError(err).Printf("Failed to read request body")
        return nil
    }
    r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

    log.WithFields(logrus.Fields{
        "body": string(body),
    }).Debug("request body")

    h := make(map[string]string)
    for k, v := range r.Header {
        h[k] = strings.Join(v, ",")
    }
    return &Record{
        ReqID:       uuid.New(),
        RequestDate: time.Now(),
        Request: RequestRecord{
            Path:    r.URL.Path,
            Method:  r.Method,
            Body:    string(body),
            Headers: Headers(h),
        },
    }
}

func (r *Record) AddResponse(response ResponseRecord) *Record {
    r.Response = response
    return r
}
