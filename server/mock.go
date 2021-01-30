package server

import (
	"fmt"
	"net/http"

	"github.com/Eldius/mock-server-go/mapper"
)

func HandleRequest(router *mapper.Router) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		router.Route(r).MakeResponse(rw)
	}
}

func StartMockServer(port int, r *mapper.Router) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleRequest(r))

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
