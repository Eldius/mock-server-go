package mapper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Eldius/mock-server-go/request"
	"gopkg.in/yaml.v3"
)

/*
Router is responsible to manage requests handling
*/
type Router struct {
	Routes   []RequestMapping
	Requests []request.Record
}

/*
NewRouter creates a new Router
*/
func NewRouter(reqMap []RequestMapping) Router {
	return Router{
		Routes:   reqMap,
		Requests: make([]request.Record, 0),
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
Add adds a new request record
*/
func (r *Router) AddRequest(req *request.Record) *Router {
	r.Requests = append(r.Requests, *req)
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

func (r *Router) Handle(rw http.ResponseWriter, req *http.Request) {
	record := request.NewRecord(req)
	r.AddRequest(record)
	mapping := r.Route(req)
	if mapping != nil {
		record.AddResponse(mapping.MakeResponse(rw))
	} else {
		rw.WriteHeader(404)
		rw.Header().Add("Content-Type", "text/plain")
		_, _ = rw.Write([]byte("Mapping not found"))
	}
}
