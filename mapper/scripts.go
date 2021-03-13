package mapper

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
	"rogchap.com/v8go"
)

const (
	engine     = "otto"
	ottoEngine = "otto"
	v8Engine   = "v8"
)

func (r *RequestMapping) parseScript(rw http.ResponseWriter, req *http.Request) (respBody string, respCode int, err error) {
	// TODO define a way to choose engine or define just one option
	switch engine {
	case ottoEngine:
		return r.parseScriptOtto(rw, req)
	case v8Engine:
		return r.parseScriptV8(rw, req)
	default:
		return r.parseScriptOtto(rw, req)
	}
}

func (r *RequestMapping) parseScriptV8(rw http.ResponseWriter, req *http.Request) (respBody string, respCode int, err error) {
	respBody = ""
	respCode = 500
	err = nil
	script := *r.Response.Script

	reqValues := map[string]interface{}{
		"body": extractBody(req.Body),
		"headers": map[string][]string{
			"content-type": {"application/json"},
		},
	}
	reqStr, _ := json.Marshal(reqValues)
	ctx, _ := v8go.NewContext() // new context with a default VM
	obj := ctx.Global()         // get the global object from the context
	if err = obj.Set("req", string(reqStr)); err != nil {
		err = fmt.Errorf("Failed to set req variable: %s\n\"%s\"\n", err.Error(), script)
		return
	}
	v, err := ctx.RunScript(script, "test.js") // executes a script on the global context
	if err != nil {
		err = fmt.Errorf(`error executing script: '%s'
script value:
"%s"`, err, script)
		return
	}
	if obj.Has("res") {
		val, _ := obj.Get("res")
		_val, _ := val.AsObject()
		if _val.Has("body") {
			tmp, _ := _val.Get("body")
			respBody = tmp.String()
		}
		if _val.Has("code") {
			tmp, _ := _val.Get("code")
			respCode = int(tmp.Integer())
		}
		log.WithFields(logrus.Fields{
			"code":  respCode,
			"body":  respBody,
			"value": v.String(),
		}).Debug("ReturningScriptValue")

		return
	}

	err = fmt.Errorf("Couldn't find 'res' variable\n\"%s\"\n", script)
	respBody = err.Error()
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
