package mapper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Eldius/mock-server-go/request"
	"github.com/robertkrimen/otto"
	//lua "github.com/yuin/gopher-lua"
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

const (
	javascriptPreffix = "script:javascript:"
)

func (r *RequestMapping) MakeResponse(rw http.ResponseWriter, req *http.Request) request.ResponseRecord {
	respRec := request.ResponseRecord{
		Headers: map[string][]string{},
	}
	for k, values := range r.Response.Headers {
		respRec.Headers[k] = append(respRec.Headers[k], values...)
		rw.Header().Add(k, strings.Join(respRec.Headers[k], "; "))
	}
	if strings.HasPrefix(r.Response.Body, javascriptPreffix) {
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
		respRec.Body = r.Response.Body
		rw.WriteHeader(r.Response.StatusCode)
		_, _ = rw.Write([]byte(r.Response.Body))
	}
	return respRec
}

func (r *RequestMapping) parseScript(rw http.ResponseWriter, req *http.Request) (respBody string, respCode int, err error) {
	vm := otto.New()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}

	request := map[string]interface{}{
		"body":    string(string(body)),
		"headers": req.Header,
	}

	_ = vm.Set("req", request)
	var value otto.Value
	script := strings.TrimPrefix(r.Response.Body, javascriptPreffix)

	if _, err = vm.Run(script); err != nil {
		err = fmt.Errorf("Failed to execute script\nerror: %v\nscript:\n%s", err, script)
		return
	} else {
		value, err = vm.Get("res")
		if err != nil {
			err = fmt.Errorf("Failed to get return variable\n%s", err.Error())
			return
		}
		obj := value.Object()
		var aux otto.Value
		aux, err = obj.Get("body")
		if err != nil {
			return
		}
		respBody = aux.String()

		aux, err = obj.Get("code")
		if err != nil {
			return
		}
		respCode, err = strconv.Atoi(aux.String())
		if err != nil {
			return
		}
		return
	}
}
