package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Eldius/mock-server-go/mapper"
)

func RouteHandler(router *mapper.Router) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var mapping mapper.RequestMapping
		err := json.NewDecoder(r.Body).Decode(&mapping)
		if err != nil {
			fmt.Println("Failed to read request body")
			fmt.Println(err.Error())
		}
		router.Add(mapping)
		rw.WriteHeader(200)
		rw.Header().Add("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(mapping)
		//_, _ = rw.Write([]byte("RouteHandler request received"))
	}
}

func StartAdminServer(port int, r *mapper.Router) {
	mux := http.NewServeMux()
	mux.HandleFunc("/route", RouteHandler(r))

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
