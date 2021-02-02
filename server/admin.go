package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Eldius/mock-server-go/mapper"
	"gopkg.in/yaml.v3"
)

func RouteHandler(router *mapper.Router) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Printf("Receiving %s request\n", r.Method)
		if r.Method == "POST" {
			var mapping mapper.RequestMapping
			err := json.NewDecoder(r.Body).Decode(&mapping)
			if err != nil {
				fmt.Println("Failed to read request body")
				fmt.Println(err.Error())
				rw.WriteHeader(500)
				_, _ = rw.Write([]byte(err.Error()))
				return
			}
			router.Add(mapping)
			rw.WriteHeader(200)
			rw.Header().Add("Content-Type", "application/json")
			_ = json.NewEncoder(rw).Encode(mapping)

		} else if r.Method == "GET" {
			rw.WriteHeader(200)
			encodeResponse(router, r, rw)
		} else {
			fmt.Println("returning: 'Method not allowed'")
			rw.WriteHeader(405)
		}
	}
}

func encodeResponse(router *mapper.Router, r *http.Request, rw http.ResponseWriter) {
	accepts := r.Header.Get("Accept")
	if strings.Contains(strings.ToLower(accepts), "application/yaml") {
		rw.Header().Add("Content-Type", "application/yaml")
		_ = yaml.NewEncoder(rw).Encode(router.Routes)
	} else {
		rw.Header().Add("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(router.Routes)
	}
}

func StartAdminServer(port int, r *mapper.Router) {
	mux := http.NewServeMux()
	mux.HandleFunc("/route", RouteHandler(r))
	mux.HandleFunc("/", AdminPanelHandler(r))

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
