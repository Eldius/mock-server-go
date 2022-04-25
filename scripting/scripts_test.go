package scripting

import (
	"fmt"
	"log"
	"testing"

	"rogchap.com/v8go"
)

func TestV8Scripting(t *testing.T) {

	script := `
        console.log(req);
        var res = {};
		res.code = 200;
		res.body = JSON.stringify({
			"contract": 123,
			"status": "OK"
		});
		res;
`

	ctx := v8go.NewContext() // new context with a default VM
	obj := ctx.Global()      // get the global object from the context
	if err := obj.Set("req", `{
		"body": {"contract": 12345, "status": "PENDING"},
		"headers": {
			"content-type": ["application/json"]
		},
	}`); err != nil {
		t.Log(err)
		t.Error(err)
	}
	v, err := ctx.RunScript(script, "test.js") // executes a script on the global context
	if err != nil {
		t.Error(err)
		t.Logf(`error executing script: '%s'
script value:
"%s"`, err, script)
		log.Printf(`error executing script: '%s'
script value:
"%s"`, err, script)
		return
	}
	if obj.Has("res") {
		val, _ := obj.Get("res")
		_val, _ := val.AsObject()
		body, _ := _val.Get("body")
		code, _ := _val.Get("code")
		fmt.Println("-", _val.Value)
		fmt.Println("-", body.String())
		fmt.Println("-", code.Integer())
		if code.Integer() != 200 {
			t.Errorf("Wrong status code: %d", code.Integer())
		}
		if body.String() != `{"contract":123,"status":"OK"}` {
			t.Errorf("Wrong body: %s", body.String())
		}
	}

	if err != nil {
		t.Errorf("Must not return error, but returned '%s'", err)
	}
	t.Log(v.String())
}
