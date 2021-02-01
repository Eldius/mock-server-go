package mapper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

/*
Router is responsible to manage requests handling
*/
type Router struct {
	Routes []RequestMapping
}

/*
NewRouter creates a new Router
*/
func NewRouter(reqMap []RequestMapping) Router {
	return Router{
		Routes: reqMap,
	}
}

/*
Add adds a new mapping
*/
func (r *Router) Add(req RequestMapping) *Router {
	r.Routes = append(r.Routes, req)
	return r
}

/*
Route returns the mapping to handle request
*/
func (r *Router) Route(req *http.Request) *RequestMapping {
	for _, m := range r.Routes {
		if compareMethod(req, m) && comparePath(req, m) {
			return &m
		}
	}
	return nil
}

func ImportMappingYaml(source string) Router {
	var r Router
	f, err := os.Open(source)
	if err != nil {
		fmt.Println("Failed to parse mapping file")
		log.Fatalln(err.Error())
	}
	defer f.Close()
	_ = yaml.NewDecoder(f).Decode(&r)
	return r
}

func compareMethod(req *http.Request, m RequestMapping) bool {
	return strings.EqualFold(strings.ToLower(m.Method), strings.ToLower(req.Method))
}

func comparePath(req *http.Request, m RequestMapping) bool {
	return strings.EqualFold(strings.ToLower(m.Path), strings.ToLower(req.URL.Path))
}
