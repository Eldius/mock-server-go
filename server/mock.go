package server

import (
	"fmt"
	"net/http"

	"github.com/Eldius/mock-server-go/mapper"
)

func HandleMockRequest(router *mapper.Router) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if router.Route(r) != nil {
			router.Route(r).MakeResponse(rw)
		} else {
			rw.WriteHeader(404)
			rw.Header().Add("Content-Type", "text/plain")
			_, _ = rw.Write([]byte("Mapping not found"))
		}
	}
}

func StartMockServer(port int, r *mapper.Router) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleMockRequest(r))

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
