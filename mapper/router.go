package mapper

import (
	"net/http"
	"os"
	"strings"

	"github.com/Eldius/mock-server-go/logger"
	"github.com/Eldius/mock-server-go/request"
	"gopkg.in/yaml.v3"
)

/*
Router is responsible to manage requests handling
*/
type Router struct {
	Routes []RequestMapping `json:"routes" yaml:"routes"`
}

var log = logger.Log()

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
Add adds a new request record
*/
func (r *Router) AddRequest(req *request.Record) *Router {
	request.Persist(req)
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
		log.WithError(err).Println("Failed to parse mapping file")
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
	mapping := r.Route(req)
	if mapping != nil {
		r.AddRequest(record.AddResponse(mapping.MakeResponse(rw, req)))
	} else {
		rw.WriteHeader(404)
		rw.Header().Add("Content-Type", "text/plain")
		record.Response = request.ResponseRecord{
			Code: 404,
			Body: "Mapping not found",
		}
		r.AddRequest(record)
		_, _ = rw.Write([]byte("Mapping not found"))
	}
}

func (r *Router) GetRequests() []request.Record {
	return request.GetRequests()
}
