package mapper

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/robertkrimen/otto"
)

func (r *RequestMapping) parseScript(rw http.ResponseWriter, req *http.Request) (respBody string, respCode int, err error) {
	return r.parseScriptOtto(rw, req)
}

func (r *RequestMapping) parseScriptV8(rw http.ResponseWriter, req *http.Request) (respBody string, respCode int, err error) {
	respBody = ""
	respCode = 200
	err = nil

	return
}

func extractBody(reader io.Reader) string {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return ""
	}
	return string(body)
}

func (r *RequestMapping) parseScriptOtto(rw http.ResponseWriter, req *http.Request) (respBody string, respCode int, err error) {
	vm := otto.New()
	body := extractBody(req.Body)

	request := map[string]interface{}{
		"body":    body,
		"headers": req.Header,
	}

	_ = vm.Set("req", request)
	var value otto.Value
	script := *r.Response.Script

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
