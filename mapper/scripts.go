package mapper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/robertkrimen/otto"
	"rogchap.com/v8go"
)

func Execute() {
	ctx, _ := v8go.NewContext()                             // creates a new V8 context with a new Isolate aka VM
	ctx.RunScript("const add = (a, b) => a + b", "math.js") // executes a script on the global context
	ctx.RunScript("const result = add(3, 4)", "main.js")    // any functions previously added to the context can be called
	val, _ := ctx.RunScript("result", "value.js")           // return a value in JavaScript back to Go
	fmt.Printf("addition result: %s", val)
}

func (r *RequestMapping) parseScript(rw http.ResponseWriter, req *http.Request) (respBody string, respCode int, err error) {
	return r.parseScriptOtto(rw, req)
}

func (r *RequestMapping) parseScriptV8(rw http.ResponseWriter, req *http.Request) (respBody string, respCode int, err error) {
	respBody = ""
	respCode = 200
	err = nil

	return
}

func (r *RequestMapping) parseScriptOtto(rw http.ResponseWriter, req *http.Request) (respBody string, respCode int, err error) {
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
